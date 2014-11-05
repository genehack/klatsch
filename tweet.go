package main

import "github.com/ChimeraCoder/anaconda"

// FIXME so yes, this is totally gross. probably what i really want
// here is a custom type with an embedded anaconda.Tweet element. Need
// to test out how that will serialize to JSON though. This works for
// the moment...

type Tweet struct {
	Contributors        []int64     `json:"contributors"`
	Coordinates         interface{} `json:"coordinates"`
	CreatedAt           string      `json:"created_at"`
	FavoriteCount       int         `json:"favorite_count"`
	Favorited           bool        `json:"favorited"`
	Id                  int64       `json:"id"`
	InReplyToScreenName string      `json:"in_reply_to_screen_name"`
	InReplyToStatusID   int64       `json:"in_reply_to_status_id"`
	InReplyToUserID     int64       `json:"in_reply_to_user_id"`
	PossiblySensitive   bool        `json:"possibly_sensitive"`
	RetweetCount        int         `json:"retweet_count"`
	Retweeted           bool        `json:"retweeted"`
	Text                string      `json:"text"`
	Truncated           bool        `json:"truncated"`
}

func (tweet *Tweet) initFromAnacondaTweet(aTweet anaconda.Tweet) {
	tweet.Contributors = aTweet.Contributors
	tweet.Coordinates = aTweet.Coordinates
	tweet.CreatedAt = aTweet.CreatedAt
	tweet.FavoriteCount = aTweet.FavoriteCount
	tweet.Favorited = aTweet.Favorited
	tweet.Id = aTweet.Id
	tweet.InReplyToScreenName = aTweet.InReplyToScreenName
	tweet.InReplyToStatusID = aTweet.InReplyToStatusID
	tweet.InReplyToUserID = aTweet.InReplyToUserID
	tweet.PossiblySensitive = aTweet.PossiblySensitive
	tweet.RetweetCount = aTweet.RetweetCount
	tweet.Retweeted = aTweet.Retweeted
	tweet.Text = aTweet.Text
	tweet.Truncated = aTweet.Truncated
}
