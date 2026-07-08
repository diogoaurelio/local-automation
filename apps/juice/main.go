package main

import (
	"fmt"
	"log"
	"os"
	"juice/clients"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func GetOldestTimestamp(hours int) int64 {
	return time.Now().Add(-time.Duration(hours) * time.Hour).Unix()
}

func main() {
	fmt.Println("Starting slack bot")

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env: %s\n", err)
		return
	}

	slackUserToken := os.Getenv("SLACK_USER_TOKEN")
	if len(slackUserToken) == 0 {
		fmt.Println("Error: no slack token provided, exiting")
		return
	}

	api := slack.New(slackUserToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	// Test auth to confirm token works
	authResp, err := api.AuthTest()
	if err != nil {
		fmt.Printf("Auth test failed: %+v\n", err)
		return
	}
	fmt.Printf("Authenticated as: %s (team: %s)\n", authResp.User, authResp.Team)

	chatService := services.NewSlackService(api)

	// List channels
	fmt.Println("Retrieving all Slack channels")

	channels, err := chatService.GetChannels()
	if err != nil {
		fmt.Printf("Error getting channels: %+v", err)
		return
	}
	oldest := GetOldestTimestamp(24)
	// list channel messages
	for _, channel := range channels {
		chatService.GetChannelMessages(&channel, oldest)
	}

}
