package config

import (
	"fmt"
	"strings"
)

// Config is a struct that contains configuration for the app
type Config struct {
	Language     string
	Topic        string
	Periodicity  int
	Cron         string
	Hashtags     string
	CacheSize    int
	ShowLanguage bool
	SafeMode     bool
	Provider     string
	Publishers   string
}

// GetHashtags return a list of hashtags from a comma separated string
func (c *Config) GetHashtags() []string {

	if c.Hashtags == "" {
		return []string{}
	}

	hs := strings.Split(c.Hashtags, ",")

	for i, h := range hs {
		hs[i] = fmt.Sprintf("#%s", strings.TrimSpace(h))
	}

	return hs
}
