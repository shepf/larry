package larry

import (
	cfg "github.com/ezeoleaf/larry/config"
	"github.com/urfave/cli/v2"
)

// GetFlags returns a list of flags used the application
func GetFlags(cfg *cfg.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "topic",
			Aliases:     []string{"t"},
			Value:       "",
			Usage:       "topic for searching repos",
			Destination: &cfg.Topic,
		},
		&cli.StringFlag{
			Name:        "lang",
			Aliases:     []string{"l"},
			Value:       "",
			Usage:       "language for searching repos",
			Destination: &cfg.Language,
		},
		&cli.IntFlag{
			Name:        "time",
			Aliases:     []string{"x"},
			Value:       15,
			Usage:       "periodicity of tweet in minutes",
			Destination: &cfg.Periodicity,
		},
		&cli.StringFlag{
			Name:        "cron",
			Aliases:     []string{"c"},
			Value:       "", //默认每天早上7点发一条 --cron="0 0 7 * * ?"
			Usage:       "cron expressions",
			Destination: &cfg.Cron,
		},
		&cli.IntFlag{
			Name:        "cache",
			Aliases:     []string{"r"},
			Value:       50,
			Usage:       "size of cache for no repeating repositories",
			Destination: &cfg.CacheSize,
		},
		&cli.StringFlag{
			Name:        "hashtag",
			Aliases:     []string{"ht"},
			Value:       "",
			Usage:       "list of comma separated hashtags",
			Destination: &cfg.Hashtags,
		},
		&cli.BoolFlag{
			Name:        "show-language",
			Aliases:     []string{"tl"},
			Value:       false,
			Usage:       "bool for allowing show the language of the repo",
			Destination: &cfg.ShowLanguage,
		},
		&cli.BoolFlag{
			Name:        "safe-mode",
			Aliases:     []string{"sf"},
			Value:       false,
			Usage:       "bool for safe mode. If safe mode is enabled, no repository is published",
			Destination: &cfg.SafeMode,
		},
		&cli.StringFlag{
			Name:        "provider",
			Aliases:     []string{"pr"},
			Value:       "github",
			Usage:       "provider where publishable content comes from",
			Destination: &cfg.Provider,
		},
		&cli.StringFlag{
			Name:        "publisher",
			Aliases:     []string{"pub"},
			Value:       "twitter",
			Usage:       "list of comma separared publishers",
			Destination: &cfg.Publishers,
		},
	}
}
