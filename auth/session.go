package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"log"
	"os"
	"sync"
	"time"
)

type Session struct {
	ExpiresAt time.Time
	ID        string
	Name      string
}

func (s Session) String() string {
	return s.ID
}

var sessions = struct {
	sync.RWMutex
	store map[string]*Session
}{
	store: make(map[string]*Session),
}

const (
	FILENAME            = "sessions.csv"
	SESSION_COOKIE_NAME = "settlersscookie"
	USER_COOKIE_NAME    = "settlersucookie"
)

var csvLocker sync.RWMutex

func createSessionId() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", nil
	}
	hash := sha256.Sum256(b)
	return hex.EncodeToString(hash[:]), nil
}

func SessionCreate(user string) (*Session, error) {
	sessionId, err := createSessionId()
	if err != nil {
		return nil, err
	}
	newSession := &Session{
		ID:        sessionId,
		Name:      user,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	sessions.Lock()
	sessions.store[sessionId] = newSession
	sessions.Unlock()

	sessionPersist(newSession)

	return newSession, nil
}

func SessionIsValid(sessionCookie string) bool {
	sessions.RLock()
	session, exists := sessions.store[sessionCookie]
	sessions.RUnlock()
	return exists && time.Now().Before(session.ExpiresAt)
}

func SessionDelete(sessionCookie string) {
	sessions.Lock()
	delete(sessions.store, sessionCookie)
	sessions.Unlock()
}

func SessionGet(sessionCookie string) *Session {
	return sessions.store[sessionCookie]
}

func SessionsLoad() {
	file, err := os.Open(FILENAME)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3
	_, err = reader.Read()
	if err != nil {
		panic(err)
	}

	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	loadedEntries := 0
	for _, row := range data {
		id := row[0]
		name := row[1]
		expirationString := row[2]
		expiration, err := time.Parse(time.RFC3339, expirationString)
		if err == nil {
			loadedEntries += 1
			sessions.store[id] = &Session{
				ID:        id,
				Name:      name,
				ExpiresAt: expiration,
			}
		}
	}
	log.Printf("Loaded %d sessions from disk", loadedEntries)
}

func sessionPersist(session *Session) {
	csvLocker.Lock()
	defer csvLocker.Unlock()
	file, err := os.OpenFile(FILENAME, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	expirationString := session.ExpiresAt.Format(time.RFC3339)
	writer.Write([]string{session.ID, session.Name, expirationString})
	writer.Flush()
}
