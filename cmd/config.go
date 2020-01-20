package cmd

import "github.com/fopina/pushit/pushit"

// Config holds the user configuration
type Config struct {
	Default  string
	Profiles map[string]pushit.Profile
}

func KeysFromMap(x map[string]pushit.Profile) (keys []string) {
	keys = make([]string, len(x))

	i := 0
	for k := range x {
		keys[i] = k
		i++
	}
	return keys
}
