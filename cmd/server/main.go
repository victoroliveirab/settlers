package main

import (
	"log"
	"net/http"

	"github.com/victoroliveirab/settlers/auth"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router"
)

func main() {
	logger.Init(true)
	router.SetupRoutes()
	auth.SessionsLoad()
	logger.Log("starting ws server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
