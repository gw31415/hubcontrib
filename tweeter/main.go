package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	//"time"

	//"github.com/dghubble/go-twitter/twitter"
	//"github.com/gw31415/hubcontrib"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
		query := fmt.Sprintf(
				`{"query": "query{viewer{contributionsCollection(from: \"%s\", to: \"%s\"){ totalCommitContributions}}}"}`, time.Now().AddDate(0, 0, -1).UTC().Format("2006-01-02T15:04:05"), time.Now().UTC().Format("2006-01-02T15:04:05"))
	res, err := oauth2.NewClient(context.Background(), src).Post(
		"https://api.github.com/graphql",
		"application/json",
		bytes.NewReader([]byte(query)),
	)
	if err != nil {
		panic(err)
	}
	var out struct {
		Data struct {
			Viewer struct {
				ContributionsCollection struct {
					TotalCommitContributions int `json:"totalCommitContributions"`
				} `json:"contributionsCollection"`
			} `json:"viewer"`
		} `json:"data"`
	}
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &out); err != nil {
		panic(err)
	}
	/*
		CONSUMER_KEY := os.Getenv("CONSUMER_KEY")
		CONSUMER_SECRET := os.Getenv("CONSUMER_SECRET")
		ACCESS_TOKEN := os.Getenv("ACCESS_TOKEN")
		ACCESS_SECRET := os.Getenv("ACCESS_SECRET")
	*/
}
