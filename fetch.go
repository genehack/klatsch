package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/spf13/cobra"
	"github.com/olivere/elastic"
)

func fetch(cmd *cobra.Command, args []string) {
	var force bool = false
	if cmd.Flag("force").Value.String() == "true" {
		force = true
	}

	exitUnlessDatabaseExists()

	db := getDatabaseHandle()
	defer db.Close()

	twitter := getTwitterApiHandle(db)

	// FIXME should also be using a 'since' value here, pulled out of the db
	v := url.Values{"count": {"200"}}
	// FIXME this wants to be GetHomeTimeline()
	timeline, err := twitter.GetUserTimeline(v)
	if err != nil {
		log.Fatal(err)
	}

	inserted, err := saveTimeline(db, timeline)
	if err != nil {
		log.Fatal(err)
	}

	if force || inserted > 0 {
		if err = writeOutTimeline(db); err != nil {
			log.Fatal(err)
		}
	}
}

func getConfig(db *sql.DB) (config map[string]string) {
	rows, err := db.Query("SELECT key,val from config")
	if err != nil {
		log.Fatal(err)
	}
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

func saveTimeline(db *sql.DB, timeline []anaconda.Tweet) (inserted int, err error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare("INSERT INTO posts (id,created,text,cruft) VALUES (?,?,?,?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	client, err := elastic.NewClient(elastic.SetURL("http://ferb:9200"))
	if err != nil {
		return 0, err
	}

	for _, tweet := range timeline {
		mahTweet := initFromAnacondaTweet(tweet)

		created, err := time.Parse(time.RubyDate, mahTweet.CreatedAt)
		if err != nil {
			return inserted, err
		}

		// TODO strip this down so cruft isn't so huge -- particularly, get rid of the 'user' field.
		cruft, err := json.Marshal(mahTweet)
		if err != nil {
			return inserted, err
		}

		if _, err = stmt.Exec(mahTweet.Id, created.Unix(), mahTweet.Text, cruft); err != nil {
			if err.Error() == "UNIQUE constraint failed: posts.id" {
				//log.Printf("Skipping tweet ID %d because already in database.\n", tweet.Id)
				continue
			} else {
				return inserted, err
			}
		}

		// Use the IndexExists service to check if a specified index exists.
		exists, err := client.IndexExists("twitter").Do()
		if err != nil {
			// Handle error
			panic(err)
		}
		if !exists {
			// Create a new index.
			createIndex, err := client.CreateIndex("twitter").Do()
			if err != nil {
				// Handle error
				panic(err)
			}
			if !createIndex.Acknowledged {
				// Not acknowledged
			}
		}

		put, err := client.Index().
			Index("twitter").
			Type("tweet").
			Id(strconv.FormatInt(mahTweet.Id, 10)).
			BodyJson(mahTweet).
			Do()
		if err != nil {
			// Handle error
			panic(err)
		}
		fmt.Printf("Indexed tweet %s to index %s, type %s\n", put.Id, put.Index, put.Type)

		inserted++
	}
	tx.Commit()

	return inserted, nil
}

func updatetime() string {
	return time.Now().Format(time.RFC850)
}

func writeOutTimeline(db *sql.DB) (err error) {
	rows, err := db.Query("SELECT id,created,text FROM posts ORDER BY id desc")
	if err != nil {
		return err
	}
	defer rows.Close()

	type KlatchTweet struct {
		Id, Timestamp, Text string
	}

	tweets := []KlatchTweet{}
	count := 0
	for rows.Next() {
		var id, text string
		var created int64
		if err := rows.Scan(&id, &created, &text); err != nil {
			return err
		}
		if matched, err := regexp.MatchString("^[RT|@]", text); err != nil {
			return err
		} else if matched {
			continue
		}

		timestamp := time.Unix(created, 0).Format(time.RFC850)
		tweets = append(tweets, KlatchTweet{id, timestamp, text})
		count = count + 1
		if count >= 20 {
			break
		}
	}

	funcMap := template.FuncMap{"updatetime": updatetime}

	tmpl := template.Must(template.New("tweets.tmpl").Funcs(funcMap).ParseFiles("tmpl/tweets.tmpl"))

	if _, err := os.Stat("root"); err != nil {
		os.Mkdir("root", 0755)
	}

	output, err := os.Create("root/timeline.html")
	if err != nil {
		return nil
	}

	return tmpl.Execute(output, tweets)

}
