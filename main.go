package main

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fopina/pushit/pushit"
	"github.com/mitchellh/go-homedir"
	toml "github.com/pelletier/go-toml"
	flag "github.com/spf13/pflag"
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
	versionPtr := flag.BoolP("verbose", "v", false, "display version")
	profilePtr := flag.StringP("profile", "p", "", "profile to use")
	outputPtr := flag.BoolP("output", "o", false, "echo input - very useful when piping commands")
	configurationPtr := flag.StringP("conf", "c", configurationFile(), "TOML configuration file")
	streamPtr := flag.BoolP("stream", "s", false, "stream the output, sending each line in separate notification")
	tailPtr := flag.IntP("lines", "l", 10, "number of lines of the input that will be pushed - ignored if --stream is used")
	helpPtr := flag.BoolP("help", "h", false, "this")

	flag.Parse()

	if *helpPtr {
		flag.Usage()
		return
	}

	if *versionPtr {
		fmt.Println("Version: " + version + " (built on " + date + ")")
		return
	}

	config, err := toml.LoadFile(*configurationPtr)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *profilePtr == "" {
		defaultProfile, ok := config.Get("default").(string)
		if !ok {
			fmt.Println("Either specify a default profile in config or use the --profile option")
			os.Exit(2)
		}
		*profilePtr = defaultProfile
	}

	var profiles Config
	config.Delete("default")
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
		scanner := bufio.NewScanner(os.Stdin)
		lineQueue := list.New()
		var l string
		for scanner.Scan() {
			l = scanner.Text()
			if *outputPtr {
				fmt.Println(l)
			}
			if *streamPtr {
				if l != "" {
					err = m(l, p.Param)
					if err != nil {
						log.Printf("ERR: %v", err)
					}
				}
			} else {
				if lineQueue.Len() < *tailPtr {
					lineQueue.PushBack(l)
				} else {
					// CIRCULate it - move Front to Back and update value
					// no new Element needs to be created
					e := lineQueue.Front()
					e.Value = l
					lineQueue.MoveToBack(e)
				}
			}
		}

		if !*streamPtr {
			var b bytes.Buffer
			for e := lineQueue.Front(); e != nil; e = e.Next() {
				b.WriteString(e.Value.(string))
				b.WriteString("\n")
			}
			err = m(b.String(), p.Param)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	} else {
		err = m(msg, p.Param)
		if err != nil {
			log.Fatal(err)
		}
	}
}
