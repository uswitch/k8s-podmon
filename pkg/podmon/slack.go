package podmon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	log "github.com/Sirupsen/logrus"
)

// SlackWebhook is the URL for the webhook
type SlackWebhook struct {
	*url.URL
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
	s := SlackWebhook{u}
	return &s, nil
}

// EventLoop for firing messages
func (webhook *SlackWebhook) EventLoop(ctx context.Context, wg *sync.WaitGroup, c chan SlackMessage) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			break
		case msg := <-c:
			resp, err := webhook.Send(msg)
			if err != nil {
				log.Errorf("Slack error: %s", err)
			} else {
				log.Debugf("Got a %d from sending the following to slack: %#v", resp, msg)
			}
		}
	}
}

// Send sends a message to a slack endpoint
func (webhook *SlackWebhook) Send(msg SlackMessage) (int, error) {
	m, err := json.Marshal(msg)
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
