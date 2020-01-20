package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/fopina/pushit/pushit"
)

func StartWeb(options *CLIOptions) {
	fmt.Println(`Up and running!

Post raw data to http://` + options.WebBind + `/, as in:

	curl http://` + options.WebBind + `/ -d 'testing 1 2 3'

This will send that data as message using the default profile.
To use a specific one, post to http://` + options.WebBind + `/PROFILE
`)
	http.HandleFunc("/", options.handler)
	log.Fatal(http.ListenAndServe(options.WebBind, nil))
}

func (c *CLIOptions) handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}
	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "body missing")
		return
	}

	p := r.URL.Path[1:]
	if p == "" {
		if p = c.Config.Default; p == "" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Either specify a default profile in config or use the --profile option")
			return
		}
	}

	profile, ok := c.Config.Profiles[p]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(
			w,
			"Profile %v not found. Available profiles: %v\n",
			p,
			strings.Join(KeysFromMap(c.Config.Profiles), ", "),
		)
		return
	}

	m := pushit.ServiceMap[profile.Service]
	if m == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Profile %v has invalid service %v", p, profile.Service)
		return
	}

	err = m(string(body), profile.Params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}
	fmt.Fprintf(w, "ok")
}
