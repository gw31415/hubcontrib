package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/gw31415/hubcontrib"
	"github.com/joho/godotenv"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/oauth2"
)

func getHubDatas(years int, months int, days int) (string, string, int) {
	hub_src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	query := fmt.Sprintf(
		`{"query": "query{viewer{contributionsCollection(from: \"%s\", to: \"%s\"){ totalCommitContributions} login }}"}`, time.Now().AddDate(-years, -months, -days).UTC().Format("2006-01-02T15:04:05"), time.Now().UTC().Format("2006-01-02T15:04:05"))
	res, err := oauth2.NewClient(context.Background(), hub_src).Post(
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
				Login string `json:"login"`
			} `json:"viewer"`
		} `json:"data"`
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &out); err != nil {
		panic(err)
	}
	svg, err := hubcontrib.Svg(out.Data.Viewer.Login)
	if err != nil {
		panic(err)
	}
	return svg, out.Data.Viewer.Login, out.Data.Viewer.ContributionsCollection.TotalCommitContributions
}

func main() {
	defer func() {
		errmsg := recover()
		if errmsg != nil {
			log.Fatal(errmsg)
			os.Exit(1)
		}
	}()

	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}
	fmt.Print("Getting GitHub Data ...")
	svg, login, commits := getHubDatas(0, 0, 1)
	fmt.Println("  Done.")

	tw_src := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	twitter := tw_src.Client(context.Background(), token)

	fmt.Print("Rasterizing the Graph ...")
	svg_data := bytes.NewReader([]byte(svg))
	icon, err := oksvg.ReadIconStream(svg_data)
	if err != nil {
		panic(err.Error())
	}
	const (
		w = 800
		h = 128
	)
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)
	b := bytes.NewBuffer([]byte{})
	err = png.Encode(b, rgba)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("  Done.")

	fmt.Println("Posting the Image ...")
	val := url.Values{}
	val.Add("media_data", base64.URLEncoding.EncodeToString(b.Bytes()))
	res, err := twitter.PostForm("https://upload.twitter.com/1.1/media/upload.json", val)
	if err != nil {
		panic(err.Error())
	}
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var upload_media struct {
		Id string `json:"media_id_string"`
	}
	if err := json.Unmarshal(out, &upload_media); err != nil {
		panic(err.Error())
	}
	fmt.Println("Tweet")
	val = url.Values{}
	val.Add("status", fmt.Sprintf("ログイン名 : %s\n本日の総コミット数 : %d", login, commits))
	val.Add("media_ids", upload_media.Id)
	res, err = twitter.PostForm("https://api.twitter.com/1.1/statuses/update.json", val)
	if err != nil {
		panic(err.Error())
	}
	out, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(out))
}
