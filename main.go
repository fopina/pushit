package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fopina/pushit/pushit"
	"github.com/mitchellh/go-homedir"
	toml "github.com/pelletier/go-toml"
)

// Config holds the user configuration
type Config map[string]pushit.Profile

func configurationFile() string {
	d, _ := homedir.Dir()
	return filepath.Join(d, ".pushit.toml")
}

var version string = "DEV"
var date string

func main() {
	versionPtr := flag.Bool("v", false, "display version")
	profilePtr := flag.String("p", "default", "profile to use")
	configurationPtr := flag.String("c", configurationFile(), "TOML configuration file")
	flag.Parse()

	if *versionPtr {
		fmt.Println("Version: " + version + " (built on " + date + ")")
		return
	}

	config, err := toml.LoadFile(*configurationPtr)
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

	p, ok := profiles[*profilePtr]
	if !ok {
		log.Fatalf("Profile %v not found in %v", *profilePtr, *configurationPtr)
	}

	m := pushit.ServiceMap[p.Service]
	if m == nil {
		log.Fatalf("Profile %v has invalid service %v", *profilePtr, p.Service)
	}

	msg := flag.Arg(0)
	if msg == "" {
		// read from STDIN (for piped commands)
		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		msg = string(input)
	}

	err = m(msg, p.Param)
	if err != nil {
		log.Fatal(err)
	}
}
