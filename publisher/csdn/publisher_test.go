package csdn

import (
	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestNewPublisher(t *testing.T) {
	res, err := http.Get("http://www.csdn.com")
	if err != nil {
		log.Fatal(err)
	}
	//利用ioutil包读取服务器返回的数据
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close() //一定要记得关闭连接
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", data)
}

func TestPublishContentInSafeMode(t *testing.T) {
	c := config.Config{SafeMode: true}
	ak := AccessKeys{}

	p := NewPublisher(ak, c)

	ti, s, u := "ti", "s", "u"

	cont := domain.Content{Title: &ti, Subtitle: &s, URL: &u}

	r, err := p.PublishContent(&cont)

	if !r {
		t.Error("expected content published in Safe Mode. No content published")
	}

	if err != nil {
		t.Errorf("expected no error got %v", err)
	}
}
