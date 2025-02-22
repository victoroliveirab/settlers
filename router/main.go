package router

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/victoroliveirab/settlers/db/models"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/types"
	wsUtils "github.com/victoroliveirab/settlers/router/ws/utils"
)

const (
	SESSION_COOKIE_NAME = "settlersscookie"
)

var upgrader = websocket.Upgrader{}

func sessionMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(SESSION_COOKIE_NAME)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		sessionID := sessionCookie.Value
		if sessionID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		session, err := models.SessionGet(db, sessionID)
		if err != nil || (session.ExpiresAt.Valid && session.ExpiresAt.Time.Before(time.Now())) {
			http.Error(w, "Session expired or invalid", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", session.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetupRoutes(db *sql.DB) {
	fs := http.FileServer(http.Dir("client"))

	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.Handle("/favicon.ico", http.StripPrefix("/", fs))

	// WS
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromCookie(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.LogError(user.ID, "upgrader.Upgrade", -1, err)
			return
		}
		defer c.Close()

		// User may have refreshed page while in a game
		for {
			_, err := wsUtils.ReadJson(c, user.ID)
			if err != nil {
				break
			}

			wsUtils.WriteJson(c, user.ID, &types.Message{
				Type: "response",
				Payload: map[string]interface{}{
					"a": 42,
				},
			})
		}
	})

	// API

	http.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		user, err := models.UserCheckCredentials(db, username, password)
		if err != nil {
			http.Error(w, "Wrong username and/or password", http.StatusBadRequest)
		}

		session, err := models.SessionCreate(db, int64(user.ID), time.Hour)
		cookie := http.Cookie{
			Name:     SESSION_COOKIE_NAME,
			Value:    session,
			MaxAge:   60 * 60,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/lobby", http.StatusSeeOther)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client/login.html")
	})

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client/signup.html")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" || r.Method != "GET" {
			w.WriteHeader(404)
			w.Write([]byte("Resource not found"))
			return
		}
		http.ServeFile(w, r, "client/index.html")
	})
}
