package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/fopina/pushit/services"
	"github.com/fopina/pushit/services/slack"
	"github.com/mitchellh/go-homedir"
	toml "github.com/pelletier/go-toml"
)

// Profile holds specific profile data (such as which Service to use and its settings)
type Profile struct {
	Service string
	Param   services.ServiceConfig
}

// Config holds the user configuration
type Config map[string]Profile

func configurationFile() string {
	d, _ := homedir.Dir()
	return filepath.Join(d, ".pushit.toml")
}

func main() {
	config, err := toml.LoadFile(configurationFile())
	if err != nil {
		fmt.Println(err)
		return
	}
	var profiles Config
	err = config.Unmarshal(&profiles)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range profiles {
		if v.Service == "slack" {
			err = slack.PushIt("Test Message from golangcode.com", v.Param)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
