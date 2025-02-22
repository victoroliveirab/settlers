package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/tursodatabase/go-libsql"
	"github.com/victoroliveirab/settlers/logger"
)

type Turso struct {
	CleanUp func()
	Db      *sql.DB
}

func TursoInit(localDb, primaryUrl, authToken string, syncInterval time.Duration) (Turso, error) {
	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		logger.LogError("system", "os.MkdirTemp", -1, err)
		return Turso{}, err
	}

	dbPath := filepath.Join(dir, localDb)
	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
		libsql.WithSyncInterval(syncInterval),
	)
	if err != nil {
		logger.LogError("system", "libsql.NewEmbeddedReplicaConnector", -1, err)
		return Turso{}, err
	}

	db := sql.OpenDB(connector)
	return Turso{
		CleanUp: func() {
			os.RemoveAll(dir)
			connector.Close()
		},
		Db: db,
	}, nil
}
