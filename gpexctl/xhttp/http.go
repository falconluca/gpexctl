package xhttp

import (
	"bytes"
	"fmt"
	log "github.com/golang/glog"
	"github.com/pkg/errors"
	"gpex/gpexctl/common"
	"gpex/gpexctl/config"
	"io"
	"net/http"
	u "net/url"
	"time"
)

type (
	ClientP struct {
		Client http.Client
	}
)

var (
	Client *ClientP
)

func InitClient() {
	var err error
	if Client, err = NewClientP(config.Conf); err != nil {
		log.Errorf("%#+v", err)
	}
}

func NewClientP(conf *config.Config) (c *ClientP, err error) {
	var httpProxy *u.URL
	if len(conf.HTTPProxy) > 0 {
		if httpProxy, err = u.Parse(fmt.Sprintf("http://%v", conf.HTTPProxy)); err != nil {
			return nil, errors.Wrap(err, "parse http proxy string failed")
		}
	}

	return &ClientP{
		Client: http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(httpProxy),
			},
			Timeout: 5 * time.Second,
		},
	}, nil
}

func (c ClientP) HandleRequest(httpMethod string, url string, reqBody []byte) []byte {
	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(reqBody))
	if err != nil {
		common.ExitWithError(err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		common.ExitWithErrorf("failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		common.ExitWithErrorf("failed: %v", err)
	}

	if !successfulStatusCode(resp.StatusCode) {
		common.ExitWithErrorf("http status code: %v, failed: %v", resp.StatusCode, err)
	}
	return body
}

func successfulStatusCode(code int) bool {
	return code >= 200 && code < 300
}
