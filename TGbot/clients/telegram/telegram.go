package telegram

import (
	"net/http"
	"path"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates() ([]Update, error) {
	u := path.Join(c.basePath, getUpdates)
	request := http.NewRequest(http.MethodGet)
}

func (c *Client) SendMessages() {

}
