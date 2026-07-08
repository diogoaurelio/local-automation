package services

import (
	"fmt"
	"time"
	"github.com/slack-go/slack"
)


type Channel struct {
	ID string
	Name string
}

type ChatService interface {
	GetChannels() ([]Channel, error)
	GetChannelMessages(channel *Channel, oldest int64) ([]string, error)
}

type SlackService struct{
	client *slack.Client
}

func NewSlackService(client *slack.Client) SlackService {
	return SlackService {client}
}

func (ss *SlackService) GetChannels() ([]Channel, error) {
	api := ss.client

	var allChannels []Channel
	cursor := ""

	for {

		params := &slack.GetConversationsParameters{
			Types: []string{"public_channel", "private_channel"},
			Limit: 200,
			Cursor: cursor,
			ExcludeArchived: true,
		}
		channels, nextCursor, err := api.GetConversations(params)
		if err != nil {
			fmt.Printf("Error getting channels: %+v\n", err)
			return nil, err
		}
		for _, ch := range channels {
			allChannels = append(allChannels, Channel{ID: ch.ID, Name: ch.Name})
		}
		if nextCursor == "" {
			break
		}
		cursor = nextCursor
		time.Sleep(2* time.Second)
	}

	return allChannels, nil
}

func (ss *SlackService) GetChannelMessages(channel *Channel, oldest int64) ([]string, error) {

	histParams := &slack.GetConversationHistoryParameters{
		ChannelID: channel.ID,
		Oldest:    fmt.Sprintf("%d", oldest),
		Limit:     200,
	}
	history, err := ss.client.GetConversationHistory(histParams)
	if err != nil {
		fmt.Printf("Error getting history for %s: %+v\n", channel.Name, err)
		return nil, err
	}
	fmt.Printf("Channel %s: %d messages in last 24h\n", channel.Name, len(history.Messages))

	return []string{""}, nil
}


// validate at compile time if interface is fully implemented
//var _ ChatService = (*ChatService)(nil)
