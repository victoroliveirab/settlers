package router

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/victoroliveirab/settlers/db/models"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/entities"

	"github.com/victoroliveirab/settlers/router/ws/handlers/connect"
	"github.com/victoroliveirab/settlers/router/ws/handlers/match"
	prematch "github.com/victoroliveirab/settlers/router/ws/handlers/pre-match"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

const (
	SESSION_COOKIE_NAME = "settlersscookie"
	USER_COOKIE_NAME    = "settlersucookie"
)

var upgrader = websocket.Upgrader{}

func SetupRoutes(db *sql.DB) {
	l := entities.NewLobby()
	fs := http.FileServer(http.Dir("client/dist"))

	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.Handle("/favicon.ico", http.StripPrefix("/", fs))

	// WS

	http.Handle("/ws", chainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("userID").(int64)
			if userID == 0 {
				// UserID must always exist here because it's called after withSessionMiddleware
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			roomID := r.URL.Query().Get("room")
			room, exists := l.GetRoom(roomID)
			if !exists {
				logger.LogMessage(userID, fmt.Sprintf("l.GetRoom(%s)", roomID), "Room doesn't exist")
				http.Error(w, "Resource not found", http.StatusNotFound)
				return
			}

			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				logger.LogError(userID, "upgrader.Upgrade", -1, err)
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			user, err := models.UserGetByID(db, userID)

			if err != nil {
				logger.LogError(userID, "models.UserGetByID", -1, err)
				http.Redirect(w, r, "/login", http.StatusUnauthorized)
				return
			}

			room.Lock()
			var existingPlayer *entities.GamePlayer
			for _, participant := range room.Participants {
				if participant.Player != nil && participant.Player.ID == userID {
					existingPlayer = participant.Player
					break
				}
			}
			room.Unlock()

			if existingPlayer != nil {
				existingPlayer.Connect(
					conn,
					func(msg *types.WebSocketClientRequest) {
						room.EnqueueIncomingMessage(existingPlayer, msg)
					},
				)
				room.ReconnectPlayer(existingPlayer)
				connect.SendCurrentGameState(existingPlayer)
			} else {
				if !room.CanAddPlayer() {
					conn.WriteJSON(map[string]string{"error": "Cannot join this room"})
					conn.Close()
					return
				}

				newPlayer := entities.NewPlayer(
					user,
					room,
					func(player *entities.GamePlayer) {
						fmt.Println(player.Username, "onDisconnect call")
					},
				)

				newPlayer.Connect(
					conn,
					func(msg *types.WebSocketClientRequest) {
						room.EnqueueIncomingMessage(newPlayer, msg)
					},
				)

				room.AddPlayer(newPlayer)

				connect.SendCurrentGameState(newPlayer)
			}
		}),
		withSessionMiddleware(db),
	))

	// API

	http.HandleFunc("POST /signup", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		name := r.FormValue("username")
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		id, err := models.UserCreate(db, username, name, email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		session, err := models.SessionCreate(db, id, username, time.Hour)
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

	http.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		userID, err := models.UserCheckCredentials(db, username, password)
		if err != nil {
			http.Error(w, "Wrong username and/or password", http.StatusBadRequest)
		}

		session, err := models.SessionCreate(db, userID, username, 30*24*time.Hour)
		cookie := http.Cookie{
			Name:     SESSION_COOKIE_NAME,
			Value:    session,
			MaxAge:   30 * 24 * 60 * 60,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}

		userCookie := http.Cookie{
			Name:   USER_COOKIE_NAME,
			Value:  username,
			MaxAge: 30 * 24 * 60 * 60,
			Path:   "/",
			Secure: true,
		}
		http.SetCookie(w, &cookie)
		http.SetCookie(w, &userCookie)
		http.Redirect(w, r, "/lobby", http.StatusSeeOther)
	})

	http.Handle("POST /create-room", chainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id := r.FormValue("id")

			room, err := l.CreateRoom(id, "base4", 4, 16)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			room.RegisterIncomingMessageHandler(prematch.TryHandle)
			room.RegisterIncomingMessageHandler(match.TryHandle)

			go room.ProcessIncomingMessages()
			go room.ProcessOutgoingMessages()
			http.Redirect(w, r, fmt.Sprintf("/game/%s", id), http.StatusSeeOther)
		}),
		withSessionMiddleware(db),
		withAdminMiddleware(db),
		withLoggingMiddleware,
	))

	// Client

	http.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client/dist/login.html")
	})

	http.HandleFunc("GET /signup", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client/dist/signup.html")
	})

	http.Handle("GET /lobby", chainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "client/dist/lobby.html")
		}),
		withSessionMiddleware(db),
	))

	http.Handle("GET /game/{id}", chainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "client/dist/game.html")
		}),
		withSessionMiddleware(db),
	))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" || r.Method != "GET" {
			w.WriteHeader(404)
			w.Write([]byte("Resource not found"))
			return
		}
		http.ServeFile(w, r, "client/dist/index.html")
	})
}
