package nwlib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SlackBody struct
type SlackBody struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	IconURL  string `json:"icon_url"`
	Channel  string `json:"channel"`
}

var cfg = LoadAWSConfig()
var snsClient = sns.New(cfg)

// PublishSNS publish message to SNS topic
func PublishSNS(arn string, message string) error {
	if message == "" {
		return errors.New("Cannot publish SNS, empty message")
	}

	fmt.Printf("Publishing message to SNS: %s\n", message)

	request := snsClient.PublishRequest(&sns.PublishInput{
		Message:  &message,
		TopicArn: &arn,
	})

	_, err := request.Send()
	if err != nil {
		fmt.Println("Problem publishing to SNS topic", err)
	}

	return err
}

// PublishSlack publish message to Slack
func PublishSlack(webhook string, message string, channel string) error {
	if message == "" {
		return errors.New("Cannot publish Slack, empty message")
	}

	fmt.Printf("Publishing message to Slack: %s\n", message)

	slackBody := &SlackBody{
		Text: message,
	}

	if channel != "" {
		slackBody.Channel = channel
	}

	payload, err := json.Marshal(slackBody)

	if err != nil {
		fmt.Println("Problem converting Slack body to json", err)
		return err
	}

	if res, err := http.Post(webhook, "application/json", bytes.NewBuffer(payload)); err != nil || res.StatusCode != http.StatusOK {
		fmt.Println("Problem with Slack", err)
		return err
	}

	return err
}