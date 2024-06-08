package testclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	apiURL     *url.URL
	httpClient *http.Client
}

type Opts struct {
	TLSCert *tls.Certificate
	URL     string
	// Middleware to wrap the default transport
	Middleware Middleware
}

type Middleware func(http.RoundTripper) http.RoundTripper

func NewClient(opts Opts) (*Client, error) {
	apiURL, err := url.Parse(opts.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid API URL %q: %w", opts.URL, err)
	}

	transport := &http.Transport{
		Proxy:        http.ProxyFromEnvironment,
		MaxIdleConns: 100,
		// .....
	}

	if opts.TLSCert != nil {
		transport.TLSClientConfig = &tls.Config{
			Certificates: []tls.Certificate{*opts.TLSCert},
		}
	}

	// http.RoundTripper is a interface. golang接口和实现之间是隐式。transport这个结构体中实现了RoundTrip这个方法
	// 并且http.RoundTripper声明了RoundTrip方法，说以可以把transport赋值给接口。
	var roundTripper http.RoundTripper = transport

	if opts.Middleware != nil {
		// 可以通过外部赋值Middleware的方式来修改http client的transport.
		roundTripper = opts.Middleware(roundTripper)
	}

	return &Client{
		apiURL:     apiURL,
		httpClient: &http.Client{Transport: roundTripper},
	}, nil
}

func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context.Background(), method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "testing")
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) Get() (string, error) {
	request, err := c.newRequest(http.MethodGet, c.apiURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("error constructing request: %w", err)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}

	log.Printf("Request made successfully %s", resp.Body)
	return "", err
}
