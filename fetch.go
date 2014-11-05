package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/spf13/cobra"
)

func fetch(cmd *cobra.Command, args []string) {
	exitUnlessDatabaseExists()

	db := getDatabaseHandle()
	defer db.Close()

	twitter := getTwitterApiHandle(db)

	v := url.Values{"count": {"200"}, "exclude_replies": {"true"}}
	timeline, err := twitter.GetUserTimeline(v)
	if err != nil {
		log.Fatal(err)
	}

	err = saveTimeline(db, timeline)
	if err != nil {
		log.Fatal(err)
	}
}

func getConfig(db *sql.DB) (config map[string]string) {
	rows, err := db.Query("SELECT key,val from config")
	defer rows.Close()

	config = make(map[string]string)

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

func saveTimeline(db *sql.DB, timeline []anaconda.Tweet) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO posts (id,created,text,cruft) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, tweet := range timeline {
		mahTweet := new(Tweet)
		mahTweet.initFromAnacondaTweet(tweet)

		created, err := time.Parse(time.RubyDate, mahTweet.CreatedAt)
		if err != nil {
			return err
		}

		// TODO strip this down so cruft isn't so huge -- particularly, get rid of the 'user' field.
		cruft, err := json.Marshal(mahTweet)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(mahTweet.Id, created.Unix(), mahTweet.Text, cruft)
		if err != nil {
			if err.Error() == "UNIQUE constraint failed: posts.id" {
				//log.Printf("Skipping tweet ID %d because already in database.\n", tweet.Id)
				continue
			}
			return err
		}
		log.Printf("Inserted tweet ID %d into database", tweet.Id)
	}
	tx.Commit()

	return nil
}
