package configuration

import (
	"database/sql"
	"errors"
	"log"

	"github.com/Bayan2019/hackathon-2025-police-volunteers/repositories/database"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// var MapCfg *MapConfiguration
var ApiCfg *ApiConfiguration

type ApiConfiguration struct {
	Conn      *sql.DB
	DB        *database.Queries
	Dir       string
	JwtSecret string
}

func Connect2DB(dbPath string) error {
	// https://github.com/libsql/libsql-client-go/#open-a-connection-to-sqld
	// libsql://[your-database].turso.io?authToken=[your-auth-token]
	if dbPath == "" {
		return errors.New("no DataBase Path")
	}
	// db, err := sql.Open("sqlite3", dbPath)
	db, err := sql.Open("libsql", dbPath)
	if err != nil {
		return err
	}

	dbQueries := database.New(db)
	ApiCfg = &ApiConfiguration{
		Conn: db,
		DB:   dbQueries,
	}
	log.Println("Connected to database!")
	return nil
}
