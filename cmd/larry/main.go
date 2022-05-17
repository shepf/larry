package main

import (
	"fmt"
	"github.com/ezeoleaf/larry/cache"
	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/larry"
	"github.com/ezeoleaf/larry/provider"
	"github.com/ezeoleaf/larry/provider/github"
	"github.com/ezeoleaf/larry/publisher"
	"github.com/ezeoleaf/larry/publisher/csdn"
	githubPub "github.com/ezeoleaf/larry/publisher/github"
	"github.com/ezeoleaf/larry/publisher/twitter"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
)

var (
	redisAddress = envString("REDIS_ADDRESS", "localhost:6379")

	//redis
	githubAccessToken      = envString("GITHUB_ACCESS_TOKEN", "")
	githubPublishRepoOwner = envString("GITHUB_PUBLISH_REPO_OWNER", "")
	githubPublishRepoName  = envString("GITHUB_PUBLISH_REPO_NAME", "")
	githubPublishRepoFile  = envString("GITHUB_PUBLISH_REPO_FILE", "README.md")

	//csdn
	CsdnCookieKey = envString("CSDN_COOKIE", "")

	//twitter
	twitterConsumerKey    = envString("TWITTER_CONSUMER_KEY", "")
	twitterConsumerSecret = envString("TWITTER_CONSUMER_SECRET", "")
	twitterAccessToken    = envString("TWITTER_ACCESS_TOKEN", "")
	twitterAccessSecret   = envString("TWITTER_ACCESS_SECRET", "")
)

func main() {
	cfg := config.Config{}

	app := &cli.App{
		Name:  "Larry",
		Usage: "Twitter bot that publishes random information from providers",
		Flags: larry.GetFlags(&cfg),
		Authors: []*cli.Author{
			{Name: "@ezeoleaf", Email: "ezeoleaf@gmail.com"},
			{Name: "@beesaferoot", Email: "hikenike6@gmail.com"}},
		Action: func(c *cli.Context) error {
			prov, err := getProvider(cfg)
			if err != nil {
				log.Fatal(err)
			}

			if prov == nil {
				log.Fatalf("could not initialize provider for %v", cfg.Provider)
			}

			pubs, err := getPublishers(cfg)
			if err != nil {
				log.Fatal(err)
			}

			if len(pubs) == 0 {
				log.Fatalln("no publishers initialized")
			}

			s := larry.Service{Provider: prov, Publishers: pubs}

			if cfg.Cron != "" {
				c := cron.New(cron.WithSeconds())
				enterId, err := c.AddFunc("0 0 7 * * ?", func() {

					err := s.Run()
					if err != nil {
						log.Printf("Error in jobTask larry.Service.Run(): %v", err)
					}
					fmt.Printf("任务启动: %s \n", time.Now().Format("2006-01-02 15:04:05"))

				})
				if err != nil {
					panic(err)
				}
				fmt.Printf("任务id是 %d \n", enterId)
				c.Start()

				select {}

			} else {
				for {
					err := s.Run()
					if err != nil {
						log.Printf("Error in larry.Service.Run(): %v", err)
					}
					time.Sleep(time.Duration(cfg.Periodicity) * time.Minute)
				}
			}

			return nil
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatalln(err)
	}
}

func getProvider(cfg config.Config) (larry.Provider, error) {
	ro := &redis.Options{
		Addr:     redisAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	cacheClient := cache.NewClient(ro)
	if cfg.Provider == provider.Github {
		np := github.NewProvider(githubAccessToken, cfg, cacheClient)
		return np, nil
	}

	return nil, nil
}

func getPublishers(cfg config.Config) (map[string]larry.Publisher, error) {

	pubs := make(map[string]larry.Publisher)

	ps := strings.Split(cfg.Publishers, ",")

	for _, v := range ps {
		v = strings.ToLower(strings.TrimSpace(v))

		if _, ok := pubs[v]; ok {
			continue
		}

		if v == publisher.Csdn {
			accessKeys := csdn.AccessKeys{
				CsdnCookie: CsdnCookieKey,
			}
			pubs[v] = csdn.NewPublisher(accessKeys, cfg)
		}

		if v == publisher.Twitter {
			accessKeys := twitter.AccessKeys{
				TwitterConsumerKey:    twitterConsumerKey,
				TwitterConsumerSecret: twitterConsumerSecret,
				TwitterAccessToken:    twitterAccessToken,
				TwitterAccessSecret:   twitterAccessSecret,
			}
			pubs[v] = twitter.NewPublisher(accessKeys, cfg)
		} else if v == publisher.Github {
			pubs[v] = githubPub.NewPublisher(githubAccessToken, cfg, githubPublishRepoOwner, githubPublishRepoName, githubPublishRepoFile)
		}

	}

	return pubs, nil
}

func envString(key string, fallback string) string {
	if value, ok := syscall.Getenv(key); ok {
		return value
	}
	return fallback
}
