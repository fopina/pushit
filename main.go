package main

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fopina/pushit/pushit"
	"github.com/mitchellh/go-homedir"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

// Config holds the user configuration
type Config struct {
	Default  string
	Profiles map[string]pushit.Profile
}

func configurationFile() string {
	d, _ := homedir.Dir()
	return filepath.Join(d, ".pushit.conf")
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

	yamlFile, err := ioutil.ReadFile(*configurationPtr)
	if err != nil {
		fmt.Println(err)
		return
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *profilePtr == "" {
		if config.Default == "" {
			fmt.Println("Either specify a default profile in config or use the --profile option")
			os.Exit(2)
		}
		*profilePtr = config.Default
	}

	p, ok := config.Profiles[*profilePtr]
	if !ok {
		log.Fatalf(
			"Profile %v not found in %v. Available profiles: %v",
			*profilePtr,
			*configurationPtr,
			strings.Join(keysFromMap(config.Profiles), ", "),
		)
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
					err = m(l, p.Params)
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
			err = m(b.String(), p.Params)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	} else {
		err = m(msg, p.Params)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func keysFromMap(x map[string]pushit.Profile) (keys []string) {
	keys = make([]string, len(x))

	i := 0
	for k := range x {
		keys[i] = k
		i++
	}
	return keys
}
