package main

import (
	"fmt"
	"html"
	"log"
	"time"

	"github.com/spf13/cobra"
)

func search(cmd *cobra.Command, args []string) {
	exitUnlessDatabaseExists()

	searchTerm := "%" + args[0] + "%"

	db := getDatabaseHandle()
	defer db.Close()

	rows, err := db.Query("SELECT created,text FROM posts WHERE text LIKE ? ORDER BY id desc", searchTerm)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var text string
		var created int64
		if err := rows.Scan(&created, &text); err != nil {
			log.Fatal(err)
		}
		timestamp := time.Unix(created, 0).Format(time.RFC850)
		fmt.Printf("-------------------\n%s\n at %s\n\n", html.UnescapeString(text), timestamp)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
