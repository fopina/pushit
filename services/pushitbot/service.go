package pushitbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fopina/pushit/services"
)

// http://fopina.github.io/tgbot-pushitbot/
type requestBody struct {
	Message string `json:"msg"`
	Format  string `json:"format"`
}

// https://github.com/fopina/tgbot-pushitbot/blob/master/pushitbot.py#L105
type replyBody struct {
	OK          bool
	Code        string
	Description string
}

// PushIt will push the message through a Slack webhook
func PushIt(msg string, config services.Config) error {
	postBody, _ := json.Marshal(requestBody{
		Message: msg,
		Format:  config["format"],
	})
	req, err := http.NewRequest(http.MethodPost, "https://tgbots.skmobi.com/pushit/"+config["token"], bytes.NewBuffer(postBody))
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
