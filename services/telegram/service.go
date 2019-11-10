package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fopina/pushit/services"
)

type requestBody struct {
	Text   string `json:"text"`
	ChatID string `json:"chat_id"`
}

type replyBody struct {
	OK          bool
	Code        int `json:"error_code"`
	Description string
}

// PushIt will push the message through a Slack webhook
func PushIt(msg string, config services.Config) error {
	postBody, _ := json.Marshal(requestBody{
		Text:   msg,
		ChatID: config["chat_id"],
	})
	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.telegram.org/bot"+config["token"]+"/sendMessage",
		bytes.NewBuffer(postBody),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	jsonReply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var reply replyBody

	err = json.Unmarshal([]byte(jsonReply), &reply)
	if err != nil {
		return err
	}

	if !reply.OK {
		return fmt.Errorf("Code %v - %v", reply.Code, reply.Description)
	}
	return nil
}
