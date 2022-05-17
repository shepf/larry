package csdn

import (
	"github.com/dghubble/sling"
	"net/http"
)

// e.g sendBlink https://blink-open-api.csdn.net/v1/pc/blink/sendBlink
const csdnBaseAPI = "https://blink-open-api.csdn.net/v1/pc/"

// Client is a Csdn client for making Csdn API requests.
type Client struct {
	sling *sling.Sling

	// Csdn API Services
	Blink *BlinkService
}

// NewClient returns a new Client.
func NewClient(httpClient *http.Client) *Client {
	base := sling.New().Client(httpClient).Base(csdnBaseAPI)
	return &Client{
		sling: base,
		Blink: newBlinkService(base.New()),
	}
}
