package pushit

import (
	"github.com/fopina/pushit/services"
	"github.com/fopina/pushit/services/slack"
	"github.com/fopina/pushit/services/pushitbot"
)

// Profile holds specific profile data (such as which Service to use and its settings)
type Profile struct {
	Service string
	Params  services.Config
}

// ServiceMap maps service name (string) to the actual PushIt method
var ServiceMap = map[string]func(string, services.Config) error{
	"slack": slack.PushIt,
	"pushitbot": pushitbot.PushIt,
}
