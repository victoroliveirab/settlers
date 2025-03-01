package router

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/victoroliveirab/settlers/db/models"
	"github.com/victoroliveirab/settlers/logger"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseRecorder) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

type Middleware func(http.Handler) http.Handler

func chainMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func withLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(recorder, r)

		duration := int(time.Since(start)) / 1e6
		logger.LogHttpRequest(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), duration, recorder.statusCode)
	})
}

func withSessionMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie(SESSION_COOKIE_NAME)
			if err != nil {
				http.Error(w, "No session found", http.StatusUnauthorized)
				return
			}
			sessionID := sessionCookie.Value
			if sessionID == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			session, err := models.SessionGet(db, sessionCookie.Value)
			if err != nil || (session.ExpiresAt.Valid && session.ExpiresAt.Time.Before(time.Now())) {
				http.Error(w, "Session expired or invalid", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", session.UserID)
			ctx = context.WithValue(ctx, "username", session.Username)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func withAdminMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("userID")
			if userID == nil {
				// UserID must always exist here because it's called after withSessionMiddleware
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			var adminValue int
			err := db.QueryRow("SELECT admin FROM Users WHERE id = ?", userID).Scan(&adminValue)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "Forbidden: admin access required", http.StatusForbidden)
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			isAdmin := adminValue == 1
			if !isAdmin {
				http.Error(w, "Forbidden: admin access required", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
