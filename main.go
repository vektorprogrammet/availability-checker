package main

import (
	"encoding/json"
	"bytes"
	"log"
	"net/http"
	"fmt"
	"time"
	"os"
)

type Configuration struct {
	SlackWorkspaces SlackWorkspaces `json:"slackWorkspaces"`
	Websites []string `json:"websites"`
}

type SlackWorkspaces []Slack

type Slack struct {
	Endpoint  string `json:"endpoint"`
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
}

type SlackMessage struct {
	Channel   string `json:"channel"`
	Text      string `json:"text"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
}

func (s Slack) Send(message string) {
	jsonData, _ := json.Marshal(SlackMessage{
		Channel: s.Channel,
		Text: message,
		Username: s.Username,
		IconEmoji: s.IconEmoji,
	})
	http.Post(s.Endpoint, "application/json", bytes.NewBuffer(jsonData))
}

var slackWorkspaces SlackWorkspaces

func checkWebsite(url string) {
	resp, err := http.Get(url)
	if err != nil {
		sendWebsiteCrashAlert(url)
		return
	}

	fmt.Println(resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		sendWebsiteCrashAlert(url)
		return
	}
}

func sendWebsiteCrashAlert(url string) {
	for _, slack := range slackWorkspaces {
		slack.Send(fmt.Sprintf(":fire::fire::fire: %s er nede!!! :fire::fire::fire:", url))
	}
}

func loadConfiguration() Configuration {
	var configuration Configuration

	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	return configuration
}

func main() {
	configuration := loadConfiguration()
	slackWorkspaces = configuration.SlackWorkspaces
	websites := configuration.Websites

	for {
		for _, website := range websites {
			checkWebsite(website)
		}
		time.Sleep(1 * time.Minute)
	}
}
