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

	"github.com/blang/semver"
	"github.com/fopina/pushit/cmd"
	"github.com/fopina/pushit/pushit"
	"github.com/mitchellh/go-homedir"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

func configurationFile() string {
	d, _ := homedir.Dir()
	return filepath.Join(d, ".pushit.conf")
}

var version string = "DEV"
var date string

const repo = "fopina/pushit"

func main() {
	options := cmd.CLIOptions{}
	versionPtr := flag.BoolP("version", "v", false, "display version")
	configurationPtr := flag.StringP("conf", "c", configurationFile(), "TOML configuration file")
	flag.StringVarP(&options.Profile, "profile", "p", "", "profile to use")
	flag.BoolVarP(&options.Output, "output", "o", false, "echo input - very useful when piping commands")
	flag.BoolVarP(&options.Stream, "stream", "s", false, "stream the output, sending each line in separate notification")
	flag.IntVarP(&options.Tail, "lines", "l", 10, "number of lines of the input that will be pushed - ignored if --stream is used")
	webPtr := flag.BoolP("web", "w", false, "Run as webserver, using raw POST data as message")
	flag.StringVarP(&options.WebBind, "web-bind", "b", "127.0.0.1:8888", "Address and port to bind web server")
	updatePtr := flag.BoolP("update", "u", false, "Auto-update pushit with latest release")
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

	if *updatePtr {
		selfUpdate()
		return
	}

	yamlFile, err := ioutil.ReadFile(*configurationPtr)
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(yamlFile, &options.Config)
	if err != nil {
		log.Fatalln(err)
	}

	if *webPtr {
		cmd.StartWeb(&options)
		return
	}

	if options.Profile == "" {
		if options.Config.Default == "" {
			log.Fatalln("Either specify a default profile in config or use the --profile option")
		}
		options.Profile = options.Config.Default
	}

	p, ok := options.Config.Profiles[options.Profile]
	if !ok {
		log.Fatalf(
			"Profile %v not found in %v. Available profiles: %v\n",
			options.Profile,
			*configurationPtr,
			strings.Join(cmd.KeysFromMap(options.Config.Profiles), ", "),
		)
	}

	m := pushit.ServiceMap[p.Service]
	if m == nil {
		log.Fatalf("Profile %v has invalid service %v", options.Profile, p.Service)
	}

	msg := flag.Arg(0)

	if msg == "" {
		// read from STDIN (for piped commands)
		scanner := bufio.NewScanner(os.Stdin)
		lineQueue := list.New()
		var l string
		for scanner.Scan() {
			l = scanner.Text()
			if options.Output {
				fmt.Println(l)
			}
			if options.Stream {
				if l != "" {
					err = m(l, p.Params)
					if err != nil {
						log.Printf("ERR: %v", err)
					}
				}
			} else {
				if lineQueue.Len() < options.Tail {
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

		if !options.Stream {
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

func selfUpdate() error {
	previous := semver.MustParse(version)
	latest, err := selfupdate.UpdateSelf(previous, repo)
	if err != nil {
		return err
	}

	if previous.Equals(latest.Version) {
		fmt.Println("Current binary is the latest version", version)
	} else {
		fmt.Println("Update successfully done to version", latest.Version)
		fmt.Println("Release note:\n", latest.ReleaseNotes)
	}
	return nil
}
