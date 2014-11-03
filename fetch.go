package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/spf13/cobra"
)

func fetch(cmd *cobra.Command, args []string) {
	exitUnlessDatabaseExists()

	db := getDatabaseHandle()
	defer db.Close()

	twitter := getTwitterApiHandle(db)

	v := url.Values{"count": {"200"}}

	// FIXME is this the right call? this command is supposed to
	// pull just your tweets, at least at the moment
	timeline, err := twitter.GetHomeTimeline(v)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(timeline) // fixme just here to quiet errors

	// FIXME stopped here
	// err = saveTimeline(db, timeline)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

func getConfig(db *sql.DB) (config map[string]string) {
	rows, err := db.Query("SELECT key,val from config")
	defer rows.Close()

	for rows.Next() {
		var key, val string
		err = rows.Scan(&key, &val)
		if err != nil {
			log.Fatal(err)
		}
		config[key] = val
	}

	err = rows.Err() // get any error encountered during iteration
	if err != nil {
		log.Fatal(err)
	}

	return
}

func getTwitterApiHandle(db *sql.DB) (api *anaconda.TwitterApi) {
	config := getConfig(db)

	anaconda.SetConsumerKey(config["consumerKey"])
	anaconda.SetConsumerSecret(config["consumerSecret"])
	api = anaconda.NewTwitterApi(config["accessToken"], config["accessSecret"])

	return
}
