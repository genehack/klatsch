package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func exitUnlessDatabaseExists() {
	if _, err := os.Stat(getDatabaseFile()); err == nil {
		return
	}

	fmt.Println("Could not find Hubonator database file.")
	fmt.Println("Have you run the `init` command?")

	os.Exit(0)
}

// TODO so this feels pretty gross...
func getDatabaseFile() string { return "hubonator.db" }

func getDatabaseHandle() *sql.DB {
	sqlfile := getDatabaseFile()

	schemaize := true
	// -f sqlfile
	if _, err := os.Stat(sqlfile); err == nil {
		schemaize = false
	}

	// TODO figure out parameterizing the database type
	db, err := sql.Open("sqlite3", sqlfile)
	if err != nil {
		log.Fatal(err)
	}

	if schemaize {
		sqlStatement := `
CREATE TABLE config (
    key TEXT PRIMARY KEY,
    val TEXT NOT NULL
);
CREATE TABLE posts (
    id TEXT PRIMARY KEY,
    created INTEGER,
    text TEXT NOT NULL,
    cruft TEXT
);
`
		_, err = db.Exec(sqlStatement)
		if err != nil {
			log.Fatal("%q: %s\n", err, sqlStatement)
		}
	}

	return db
}
