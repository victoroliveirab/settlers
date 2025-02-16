package router

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/victoroliveirab/settlers/auth"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/types"
	wsUtils "github.com/victoroliveirab/settlers/router/ws/utils"
)

var upgrader = websocket.Upgrader{}

func SetupRoutes() {
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

		name := r.FormValue("name")
		session, err := auth.SessionCreate(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		sessionCookie := http.Cookie{
			Name:     auth.SESSION_COOKIE_NAME,
			Value:    session.ID,
			MaxAge:   200 * 60 * 60 * 24 * 30,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		}

		userCookie := http.Cookie{
			Name:   auth.USER_COOKIE_NAME,
			Value:  session.Name,
			MaxAge: 200 * 60 * 60 * 24 * 30,
			Path:   "/",
			Secure: true,
		}

		http.SetCookie(w, &sessionCookie)
		http.SetCookie(w, &userCookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		_, err := getUserFromCookie(r)
		if err == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return

		}
		http.ServeFile(w, r, "client/login.html")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" || r.Method != "GET" {
			w.WriteHeader(404)
			w.Write([]byte("Resource not found"))
			return
		}
		_, err := getUserFromCookie(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.ServeFile(w, r, "client/index.html")
	})
}
