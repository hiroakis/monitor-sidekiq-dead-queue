package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Payload struct {
	Text       string `json:"text"`
	Username   string `json:"username"`
	Icon_emoji string `json:"icon_emoji"`
	Icon_url   string `json:"icon_url"`
	Channel    string `json:"channel"`
}

func postSlack(incomingUrl, channel, username, iconEmoji, iconUrl, text string) {
	emojiString := fmt.Sprintf(":%s:", iconEmoji)
	channelString := fmt.Sprintf("#%s", channel)
	quotedText := fmt.Sprintf("```\n%s\n```", text)

	payload := Payload{
		Text:       quotedText,
		Username:   username,
		Icon_emoji: emojiString,
		Icon_url:   iconUrl,
		Channel:    channelString,
	}
	params, _ := json.Marshal(payload)

	resp, _ := http.PostForm(
		incomingUrl,
		url.Values{"payload": {string(params)}},
	)
	defer resp.Body.Close()
}

func getLastPosition(filepath string) string {
	var t string

	pos, err := ioutil.ReadFile(filepath)
	if err != nil {
		t = fmt.Sprintf("%f", float64(time.Now().Unix()))
	} else {
		f, _ := strconv.ParseFloat(strings.Trim(string(pos), "\n"), 64)
		t = fmt.Sprintf("%f", f)
	}
	return t
}

func writeLastPosition(filepath, lastPosition string) {
	b := []byte(lastPosition)
	ioutil.WriteFile(filepath, b, os.FileMode(int(0644)))
}

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "6379", "Port")
	optPosfile := flag.String("posfile", "/tmp/last.txt", "Temp file name")
	optDeadQueue := flag.String("deadqueue", "sidekiq:dead", "The name of sidekiq dead queue")
	optIncomingUrl := flag.String("url", "", "slack webhook url")
	optSlackChannel := flag.String("channel", "", "slack channel")
	optSlackUsername := flag.String("slackusername", "", "slack username")
	optIconEmoji := flag.String("iconemoji", "", "slack icon emoji")
	optIconUrl := flag.String("iconurl", "", "slack icon url")
	flag.Parse()

	target := fmt.Sprintf("%s:%s", *optHost, *optPort)
	con, err := redis.Dial("tcp", target)
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()

	pos := getLastPosition(*optPosfile)
	p := "(" + pos
	// zrangebyscore 'sidekiq:dead' "(1441721371.6570976" +inf withscores
	result, err := redis.Strings(con.Do("zrangebyscore", *optDeadQueue, p, "+inf", "withscores"))
	if err != nil {
		fmt.Println(err)
	}

	if len(result) == 0 {
		writeLastPosition(*optPosfile, pos)
	} else {
		for i, text := range result {
			if i%2 == 0 {
				postSlack(*optIncomingUrl, *optSlackChannel, *optSlackUsername, *optIconEmoji, *optIconUrl, text)
			}
		}
		lastScore := result[len(result)-1]
		writeLastPosition(*optPosfile, lastScore)
	}
}
