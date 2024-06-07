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

type Config struct {
	TLSCert *tls.Certificate
	URL     string
}

func (c *Client) SetHttpClient(httpClient *http.Client) {
	c.httpClient = httpClient
}

func NewTestClient(cfg Config) (*Client, error) {
	apiURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid API URL %q: %w", cfg.URL, err)
	}

	return &Client{
		apiURL:     apiURL,
		httpClient: http.DefaultClient,
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
