package models

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"sync"
	"time"
)

type Session struct {
	SessionID string
	UserID    int64
	CreatedAt time.Time
	ExpiresAt sql.NullTime
}

var localCache = struct {
	sync.RWMutex
	store map[string]Session
}{
	store: make(map[string]Session),
}

func cacheSet(session Session) {
	localCache.Lock()
	defer localCache.Unlock()
	localCache.store[session.SessionID] = session
}

func cacheGet(sessionID string) (Session, bool) {
	localCache.RLock()
	defer localCache.RUnlock()
	session, exists := localCache.store[sessionID]
	return session, exists
}

func cacheDelete(sessionID string) {
	localCache.Lock()
	defer localCache.Unlock()
	delete(localCache.store, sessionID)
}

func createSessionId() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", nil
	}
	hash := sha256.Sum256(b)
	return hex.EncodeToString(hash[:]), nil
}

func SessionCreate(db *sql.DB, userID int64, duration time.Duration) (string, error) {
	sessionID, err := createSessionId()
	if err != nil {
		return "", err
	}
	expiresAt := time.Now().Add(duration)
	_, err = db.Exec(`INSERT INTO Sessions (session_id, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)`,
		sessionID, userID, time.Now(), expiresAt)
	if err != nil {
		return "", err
	}

	session := Session{SessionID: sessionID, UserID: userID, CreatedAt: time.Now(), ExpiresAt: sql.NullTime{Time: expiresAt, Valid: true}}
	cacheSet(session)

	return sessionID, nil
}

func SessionGet(db *sql.DB, sessionID string) (*Session, error) {
	if session, exists := cacheGet(sessionID); exists {
		return &session, nil
	}

	var session Session
	var expiresAt sql.NullTime
	err := db.QueryRow(`SELECT session_id, user_id, created_at, expires_at FROM sessions WHERE session_id = ?`, sessionID).
		Scan(&session.SessionID, &session.UserID, &session.CreatedAt, &expiresAt)
	if err != nil {
		return nil, err
	}

	session.ExpiresAt = expiresAt
	cacheSet(session)
	return &session, nil
}

func SessionDelete(db *sql.DB, sessionID string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE session_id = ?`, sessionID)
	if err != nil {
		return err
	}

	cacheDelete(sessionID)
	return nil
}
