package ysx

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	DefaultURL                  = "https://meetingpre.125339.ebupt.net"
	DefaultTokenRefreshInterval = 6 * time.Hour
	DefaultRetryLimit           = 3
	DefaultRetryInterval        = 1 * time.Second
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client interface {
	token
	meeting
}

type client struct {
	c Doer
	r *redis.Client

	mu                   sync.RWMutex
	url                  string
	basicIdentity        string
	basicMobile          string
	basicKey             string
	tokenRefreshInterval time.Duration
	retryLimit           int
	retryInterval        time.Duration
	header               http.Header
	option               *redis.Options
	tokenKey             string
}

type ClientOptionFunc func(*client) error

func SetBasicAuth(identity, mobile, key string) ClientOptionFunc {
	return func(c *client) error {
		c.basicIdentity = identity
		c.basicMobile = mobile
		c.basicKey = key

		return nil
	}
}

func SetURL(url string) ClientOptionFunc {
	return func(c *client) error {
		c.url = url

		return nil
	}
}

func SetTokenRefreshInterval(interval time.Duration) ClientOptionFunc {
	return func(c *client) error {
		c.tokenRefreshInterval = interval

		return nil
	}
}

func SetRetryLimit(limit int) ClientOptionFunc {
	return func(c *client) error {
		c.retryLimit = limit

		return nil
	}
}

func SetRetryInterval(interval time.Duration) ClientOptionFunc {
	return func(c *client) error {
		c.retryInterval = interval

		return nil
	}
}

func SetProxyURL(u string) ClientOptionFunc {
	if u == "" {
		return func(c *client) error {
			return nil
		}
	}

	return func(c *client) error {
		if client, ok := c.c.(*http.Client); ok {
			if t, ok := client.Transport.(*http.Transport); ok {
				t.Proxy = func(request *http.Request) (*url.URL, error) {
					return url.Parse(u)
				}
			}
		}

		return nil
	}
}

func SetRedisOption(option *redis.Options) ClientOptionFunc {
	return func(c *client) error {
		c.option = option

		return nil
	}
}

func SetTokenKey(key string) ClientOptionFunc {
	return func(c *client) error {
		c.tokenKey = key

		return nil
	}
}

func NewClient(options ...ClientOptionFunc) (Client, error) {
	client := &client{
		c: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		url:                  DefaultURL,
		tokenRefreshInterval: DefaultTokenRefreshInterval,
		retryLimit:           DefaultRetryLimit,
		retryInterval:        DefaultRetryInterval,
		header:               make(map[string][]string, 2),
	}

	for _, option := range options {
		if err := option(client); nil != err {
			return nil, err
		}
	}
	client.header.Add("Content-Type", "application/json; charset=utf-8")

	client.r = redis.NewClient(client.option)
	if _, err := client.r.Ping(context.Background()).Result(); nil != err {
		panic(err)
	}

	_, err := client.GetToken()
	if err != nil {
		panic(err)
	}

	go refreshToken(client)

	return client, nil
}

type PerformRequestOptions struct {
	Method  string
	Path    string
	Params  url.Values
	Body    interface{}
	Headers http.Header
}

func (c *client) performRequest(option PerformRequestOptions) (*Response, error) {
	c.mu.RLock()
	defaultHeaders := c.header
	c.mu.RUnlock()

	var (
		err error
		req *Request
		rsp *Response
	)

	pathWithParams := option.Path
	if len(option.Params) > 0 {
		pathWithParams += "?" + option.Params.Encode()
	}

	req, err = NewRequest(option.Method, c.url+pathWithParams)
	if err != nil {
		return nil, err
	}
	for key, value := range defaultHeaders {
		for _, v := range value {
			req.Header.Add(key, v)
		}
	}
	for key, value := range option.Headers {
		for _, v := range value {
			req.Header.Add(key, v)
		}
	}
	if nil != option.Body {
		err = req.SetBody(option.Body)
		if err != nil {
			return nil, err
		}
	}

	for n := 0; n < c.retryLimit-1; n++ {
		rsp, err = c.request((*http.Request)(req))
		if err == nil {
			break
		}

		time.Sleep(c.retryInterval)
	}
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *client) request(req *http.Request) (*Response, error) {
	r, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	rsp, err := c.newResponse(r)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}
