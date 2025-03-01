package main

import (
	"log"
	"net/http"
	"os"
	"time"

	mapsdefinitions "github.com/victoroliveirab/settlers/core/maps"
	"github.com/victoroliveirab/settlers/db"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router"
)

func main() {
	logger.Init(true)
	tursoURL := os.Getenv("TURSO_URL")
	tursoToken := os.Getenv("TURSO_TOKEN")
	turso, err := db.TursoInit("local.db", tursoURL, tursoToken, 10*time.Minute)
	if err != nil {
		panic(err)
	}
	db := turso.Db
	defer turso.CleanUp()
	mapsdefinitions.LoadMap("base4")
	router.SetupRoutes(db)
	logger.Log("starting ws server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
