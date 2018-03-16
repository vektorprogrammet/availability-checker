package main

import (
	"encoding/json"
	"bytes"
	"log"
	"net/http"
	"fmt"
	"time"
	"os"
	"sync"
)

type Configuration struct {
	SlackWorkspaces SlackWorkspaces `json:"slackWorkspaces"`
	Websites        []string        `json:"websites"`
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
		Channel:   s.Channel,
		Text:      message,
		Username:  s.Username,
		IconEmoji: s.IconEmoji,
	})
	http.Post(s.Endpoint, "application/json", bytes.NewBuffer(jsonData))
}

var slackWorkspaces SlackWorkspaces

func checkWebsite(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}

func sendWebsiteCrashAlert(url string) {
	for _, slack := range slackWorkspaces {
		slack.Send(fmt.Sprintf(":fire::fire::fire: %s is offline! :fire::fire::fire:", url))
	}
}

func sendWebsiteBackUpAlert(url string) {
	for _, slack := range slackWorkspaces {
		slack.Send(fmt.Sprintf(":male-firefighter: %s is back online! :female-firefighter:", url))
	}
}

func monitorWebsite(website string) {
	isRunning := true
	for {
		wasRunning := isRunning
		isRunning = checkWebsite(website)
		if wasRunning && !isRunning {
			sendWebsiteCrashAlert(website)
		} else if !wasRunning && isRunning {
			sendWebsiteBackUpAlert(website)
		}
		time.Sleep(1 * time.Minute)
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

	var wg sync.WaitGroup
	wg.Add(1)

	for _, website := range websites {
		go monitorWebsite(website)
	}

	wg.Wait()
}
