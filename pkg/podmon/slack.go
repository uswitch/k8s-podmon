package podmon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SlackWebhook is the URL for the webhook
type SlackWebhook struct {
	URL *url.URL
}

// SlackMessage Simple message struct for sending messages to slack
type SlackMessage struct {
	Username string `json:"username"`
	Icon     string `json:"icon_emoji"`
	Text     string `json:"text"`
	Channel  string `json:"channel"`
}

// NewSlackEndpoint is a slack webhook
func NewSlackEndpoint(slack string) (*SlackWebhook, error) {
	u, err := url.Parse(slack)
	if err != nil {
		return nil, fmt.Errorf("Error parsing URL for Slack Webhook: %s", err)
	}
	s := SlackWebhook{
		URL: u,
	}
	return &s, nil
}

// Send sends a message to a slack endpoint
func (webhook *SlackWebhook) Send(message SlackMessage) (int, error) {
	m, err := json.Marshal(message)
	if err != nil {
		return 0, err
	}
	u := webhook.URL.String()
	buf := bytes.NewReader(m)
	r, err := http.Post(u, "application/json", buf)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()
	return r.StatusCode, nil
}
