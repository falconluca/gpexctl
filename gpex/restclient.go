package gpex

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	u "net/url"
)

type RestClient struct {
	Client http.Client
}

func NewRestClient(conf *Config) (r *RestClient, err error) {
	var httpProxy *u.URL
	if len(conf.HttpProxy) > 0 {
		if httpProxy, err = u.Parse(fmt.Sprintf("http://%v", conf.HttpProxy)); err != nil {
			return nil, errors.Wrap(err, "parse http proxy string failed")
		}
	}

	return &RestClient{
		Client: http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(httpProxy),
			},
		},
	}, nil
}

func (c RestClient) Do(req *http.Request) (body []byte, err error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http request failed")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, errors.Wrap(err, "read http response body failed")
	}
	return body, nil
}
