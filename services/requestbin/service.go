package requestbin

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
	Text   string            `json:"text"`
	Params map[string]string `json:"params"`
}

type replyBody struct {
	Success bool
}

// PushIt will push the message through a Slack webhook
func PushIt(msg string, config services.Config) error {
	postBody, _ := json.Marshal(requestBody{
		Text:   msg,
		Params: config,
	})
	req, err := http.NewRequest(
		http.MethodPost,
		config["url"],
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

	if !reply.Success {
		return fmt.Errorf(string(jsonReply))
	}
	return nil
}
