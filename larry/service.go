package larry

import (
	"log"
)

// Service represents the application struct
type Service struct {
	Publishers map[string]Publisher
	Provider   Provider
	Logger     log.Logger
}

// Run executes the application
func (s Service) Run() error {
	//从仓库中取中自己要的内容到 Content 结构
	content, err := s.Provider.GetContentToPublish()
	if err != nil {
		return err
	}
	log.Println(*content.URL)

	for _, pub := range s.Publishers {
		pub.PublishContent(content)
	}

	return nil
}
