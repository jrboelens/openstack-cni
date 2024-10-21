package cniclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/containernetworking/cni/pkg/types"
	"github.com/jboelensns/openstack-cni/pkg/util"
	"github.com/joho/godotenv"
)

func New(opts *ClientOpts) (*Client, error) {
	if opts != nil {
		return &Client{Opts: *opts}, nil
	}

	// attempt to read config file
	configFile := util.Getenv("CNI_CONFIG_FILE", "/etc/cni/net.d/openstack-cni.conf")
	exists, err := util.FileExists(configFile)
	if err != nil {
		return nil, err
	}
	if exists {
		if err := godotenv.Load(configFile); err != nil {
			return nil, err
		}
	}

	timeout, err := time.ParseDuration(fmt.Sprintf("%ss", util.Getenv("CNI_REQUEST_TIMEOUT", "60")))
	if err != nil {
		return nil, err
	}

	return &Client{
		Opts: ClientOpts{
			BaseUrl:        util.Getenv("CNI_API_URL", "http://127.0.0.1:4242"),
			RequestTimeout: timeout,
			LogFileName:    util.Getenv("CNI_LOG_FILE_NAME", ""),
		},
	}, nil
}

type Client struct {
	Opts ClientOpts
}

func (me *Client) Url(path string) string {
	return fmt.Sprintf("%s%s", me.Opts.BaseUrl, path)
}

func (me *Client) CniCommand(cmd util.CniCommand) (*http.Response, error) {
	url := me.Url("/cni")

	body, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	return me.Post(url, body)
}

func (me *Client) Get(url string) (*http.Response, error) {
	return me.doRequest(url, http.MethodGet, nil)
}

func (me *Client) Post(url string, body []byte) (*http.Response, error) {
	return me.doRequest(url, http.MethodPost, &body)
}

func (me *Client) Delete(url string) (*http.Response, error) {
	return me.doRequest(url, http.MethodDelete, nil)
}

func (me *Client) doRequest(url string, method string, body *[]byte) (*http.Response, error) {
	// prepare the request with a deadline
	deadline := time.Now().Add(me.Opts.RequestTimeout)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	var req *http.Request
	var err error
	if method == http.MethodPost {
		//TODO: <.> Make this work without using strings
		bodyReader := strings.NewReader(string(*body))
		req, err = http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return nil, err
		}
		req.Header.Add("content-type", "application/json")
	} else if method == http.MethodGet || method == http.MethodDelete {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}

	// send the request
	return http.DefaultClient.Do(req)
}

func (me *Client) HandleResponse(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return []byte{}, err
	}
	return me.handleResponse(resp)
}

func (me *Client) handleResponse(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusInternalServerError {
		var e types.Error
		if err := util.FromJson(body, &e); err != nil {
			return body, err
		}
		return body, &e
	} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return body, fmt.Errorf("received invalid response %d", resp.StatusCode)
	}
	return body, nil
}

type ClientOpts struct {
	BaseUrl        string
	RequestTimeout time.Duration
	LogFileName    string
}
