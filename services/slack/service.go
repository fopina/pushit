package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/fopina/pushit/services"
)

// https://api.slack.com/custom-integrations/incoming-webhooks
type slackRequestBody struct {
	Text      string `json:"text"`
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	IconURL   string `json:"icon_url"`
}

// PushIt will push the message through a Slack webhook
func PushIt(msg string, config services.ServiceConfig) error {
	slackBody, _ := json.Marshal(slackRequestBody{
		Text:      msg,
		Channel:   config["channel"],
		Username:  config["username"],
		IconEmoji: config["icon_emoji"],
		IconURL:   config["icon_url"],
	})
	req, err := http.NewRequest(http.MethodPost, config["url"], bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack: " + buf.String())
	}
	return nil
}
