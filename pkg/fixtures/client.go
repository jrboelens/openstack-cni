package fixtures

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

type Client struct {
}

type DoOpts struct {
	ContentType    string
	RequestTimeout time.Duration
}

func DefaultDoOpts() *DoOpts {
	return &DoOpts{
		ContentType:    "application/json",
		RequestTimeout: time.Second * 5,
	}
}

func (me *Client) Get(url string, opts *DoOpts) (*http.Response, error) {
	return me.Do(http.MethodGet, url, nil, opts)
}

func (me *Client) Post(url string, body []byte, opts *DoOpts) (*http.Response, error) {
	bodyReader := bytes.NewBuffer(body)
	return me.Do(http.MethodPost, url, bodyReader, opts)
}

func (me *Client) Do(method, url string, body io.Reader, opts *DoOpts) (*http.Response, error) {
	if opts == nil {
		opts = DefaultDoOpts()
	}

	deadline := time.Now().Add(opts.RequestTimeout)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	if method == http.MethodPost {
		req.Header.Add("content-type", opts.ContentType)
	}

	// send the request
	return http.DefaultClient.Do(req)
}
