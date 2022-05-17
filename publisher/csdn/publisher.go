package csdn

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
)

// Publisher represents the publisher client
type Publisher struct {
	Client *Client
	Config config.Config
}

// AccessKeys represents the keys and tokens or cookcies needed for comunication with the client
type AccessKeys struct {
	CsdnCookie string
}

// NewPublisher returns a new publisher
func NewPublisher(accessKeys AccessKeys, cfg config.Config) Publisher {
	// accessKeys 暂时不用，直接在发送时，我们从文件读物csdn cookie
	log.Print("New Csdn Publisher")

	client := NewClient(&http.Client{})

	p := Publisher{
		Config: cfg,
		Client: client,
	}

	return p
}

// 规整要发布的内容
// prepareTweet convers a domain.Content in a string Csdn
func (p Publisher) prepareCsdn(content *domain.Content) string {
	//标题
	csdn := "👁️‍「Go仓库分享时间」‍\n"
	csdn += fmt.Sprintf("👉 仓库名: %s  \n", *content.Title)

	csdn += "👇 仓库描述: \n"
	csdn += fmt.Sprintf("%s\n", *content.Subtitle)

	csdn += fmt.Sprintf("%s \n", strings.Join(content.ExtraData, "	"))

	csdn += fmt.Sprintf("👉 仓库地址: %s \n", *content.URL)

	log.Println("prepareCsdn: " + csdn)
	return csdn
}

// PublishContent receives a content to publish and try to publish
func (p Publisher) PublishContent(content *domain.Content) (bool, error) {

	//准备要发布的内容
	csdn := p.prepareCsdn(content)
	log.Print("prepareCsdn: \n ", csdn)

	if p.Config.SafeMode {
		log.Print("Running in Safe Mode")
		log.Print(csdn)
		return true, nil
	}

	_, _, err := p.Client.Blink.Update(csdn, nil)

	if err != nil {
		log.Print(err)
		return false, err
	}

	log.Println("Content Published")
	return true, nil
}
