package github

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/ezeoleaf/larry/cache"
	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
	"github.com/go-redis/redis/v8"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type searchClient interface {
	Repositories(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error)
}
type userClient interface {
	Get(ctx context.Context, user string) (*github.User, *github.Response, error)
}

// Provider represents the provider client
type Provider struct {
	GithubSearchClient searchClient
	GithubUserClient   userClient
	CacheClient        cache.Client
	Config             config.Config
}

const emptyChar = " "

// NewProvider returns a new provider client
func NewProvider(apiKey string, cfg config.Config, cacheClient cache.Client) Provider {
	log.Print("New Github Provider")
	p := Provider{Config: cfg, CacheClient: cacheClient}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiKey},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	p.GithubSearchClient = github.NewClient(tc).Search
	p.GithubUserClient = github.NewClient(tc).Users

	return p
}

// GetContentToPublish returns a string with the content to publish to be used by the publishers
func (p Provider) GetContentToPublish() (*domain.Content, error) {
	r, err := p.getRepo()
	if err != nil {
		return nil, err
	}

	if r == nil {
		log.Printf("GetContentToPublish fail，")
		return nil, err
	}

	return p.getContent(r), nil

}

func (p Provider) getRepositories(randomChar string) ([]*github.Repository, int, error) {
	so := github.SearchOptions{ListOptions: github.ListOptions{PerPage: 1}, TextMatch: true}

	_, t, e := p.GithubSearchClient.Repositories(context.Background(), p.getQueryString(randomChar), &so)

	if e != nil {
		return nil, -1, e
	}

	return nil, t.LastPage, nil
}

//if rand.Intn(11) > 2 {
//	return emptyChar
//}
func (p Provider) getRandomChar() string {

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	return string(letters[rand.Intn(len(letters))])
}

func (p Provider) getRepo() (*github.Repository, error) {
	log.Printf("getRepo start，开始获取仓库")

	rand.Seed(time.Now().UTC().UnixNano())
	rc := p.getRandomChar()

	_, total, err := p.getRepositories(rc)
	if err != nil {
		return nil, err
	}

	if total < 1 {
		log.Printf("char %s returned 0 repositories\n", rc)
		return nil, fmt.Errorf("char %s returned 0 repositories", rc)
	}

	var repo *github.Repository

	var found bool
	var count int64 = 0

	for !found {
		randPos := rand.Intn(total)

		repo = p.getSpecificRepo(rc, randPos)

		found = repo != nil && p.isRepoNotInCache(*repo.ID)

		if found && *repo.Archived {
			log.Printf("repository archived, %s", *repo.ID)
			found = false
			log.Print("repository archived")
			log.Print(*repo.ID)
		}

		//我们只要 star 大于等于 x 的 仓库
		if found && *repo.StargazersCount < 10 {
			count += 1
			log.Printf("repository star < 10, %s ,%d", *repo.ID, count)
			found = false
			log.Print("repository star < 10")
			log.Print(*repo.ID)
		}

	}

	return repo, nil
}

func (p Provider) getQueryString(randomChar string) string {
	var qs string

	if p.Config.Topic != "" && p.Config.Language != "" {
		qs = fmt.Sprintf("topic:%s+language:%s", p.Config.Topic, p.Config.Language)
	} else if p.Config.Topic != "" {
		qs = fmt.Sprintf("topic:%s", p.Config.Topic)
	} else {
		qs = fmt.Sprintf("language:%s", p.Config.Language)
	}

	if randomChar != emptyChar {
		qs = fmt.Sprintf("%s+%s", randomChar, qs)
	}

	return qs
}

func (p Provider) getSpecificRepo(randomChar string, pos int) *github.Repository {
	so := github.SearchOptions{ListOptions: github.ListOptions{PerPage: 1, Page: pos}, TextMatch: true}

	repositories, _, e := p.GithubSearchClient.Repositories(context.Background(), p.getQueryString(randomChar), &so)

	if e != nil {
		return nil
	}

	if len(repositories.Repositories) == 0 {
		println(`这是个空数组`)
		//log.Fatalf("repositories:  %v", repositories)
		return nil
	}

	return repositories.Repositories[0]

}

func (p Provider) isRepoNotInCache(repoID int64) bool {
	k := p.Config.Topic + "-" + strconv.FormatInt(repoID, 10)
	_, err := p.CacheClient.Get(k)

	switch {
	case err == redis.Nil:
		err := p.CacheClient.Set(k, true, time.Duration(p.Config.Periodicity)*time.Minute)
		if err != nil {
			return false
		}

		return true
	case err != nil:
		log.Println("Get failed", err)
		return false
	}

	return false
}

//从仓库拿取指定信息到 自定义 Content结构
func (p Provider) getContent(repo *github.Repository) *domain.Content {
	log.Println("getContent start: repo.Name: ", repo.Name)

	c := domain.Content{Title: repo.Name, Subtitle: repo.Description, URL: repo.HTMLURL, ExtraData: []string{}}

	if p.Config.ShowLanguage && repo.Language != nil {
		l := "项目语言: " + *repo.Language
		log.Println(l)
		c.ExtraData = append(c.ExtraData, l)
	}

	if repo.StargazersCount != nil {
		stargazers := "⭐️star数:" + strconv.Itoa(*repo.StargazersCount)
		c.ExtraData = append(c.ExtraData, stargazers)
	}

	if repo.ForksCount != nil {
		forks := "✨️Fork数:" + strconv.Itoa(*repo.ForksCount)
		c.ExtraData = append(c.ExtraData, forks)
	}

	owner := p.getRepoUser(repo.Owner)
	if owner != "" {
		author := "Author: @" + owner
		c.ExtraData = append(c.ExtraData, author)
	}

	hs := p.Config.GetHashtags()
	hashtags := ""

	if len(hs) == 0 {
		if p.Config.Topic != "" {
			hashtags += "#" + p.Config.Topic + " "
		} else if p.Config.Language != "" {
			hashtags += "#" + p.Config.Language + " "
		} else if repo.Language != nil {
			hashtags += "#" + *repo.Language + " "
		}
	} else {
		for _, h := range hs {
			if hashtags != "" {
				hashtags += " "
			}
			hashtags += h
		}
	}

	c.ExtraData = append(c.ExtraData, hashtags)

	return &c
}

func (p Provider) getRepoUser(owner *github.User) string {
	if owner == nil || owner.Login == nil {
		return ""
	}

	gUser, _, err := p.GithubUserClient.Get(context.Background(), *owner.Login)

	if err != nil {
		return ""
	}

	return gUser.GetTwitterUsername()
}
