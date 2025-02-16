package router

import (
	"errors"
	"net/http"

	"github.com/victoroliveirab/settlers/auth"
)

func getUserFromCookie(r *http.Request) (*auth.Session, error) {
	sessionCookie, err := r.Cookie(auth.SESSION_COOKIE_NAME)
	if err != nil {
		return nil, errors.New("No session cookie")
	}

	cookie := sessionCookie.Value
	if !auth.SessionIsValid(cookie) {
		return nil, errors.New("Invalid session cookie")
	}

	user := auth.SessionGet(cookie)
	return user, nil
}
