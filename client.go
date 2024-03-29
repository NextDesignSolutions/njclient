package njclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "0.1.0"
	userAgent      = "nextjtag_client/ + libraryVersion"
	defaultBaseURL = "http://192.168.2.57:12345"
	mediaType      = "application/json"
	format         = "json"
)

type VersionInfo struct {
	Minor   int    `json:"minor"`
	Major   int    `json:"major"`
	Sha1    string `json:"sha1"`
	Version string `json:"version"`
}

type Config struct {
	APIVersion string
}

type NjtError struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type MinValueMaxError struct {
	Min   float32   `json:"min`
	Value float32   `json:"max`
	Max   float32   `json:"max`
	Error *NjtError `json:"error"`
}

type Client struct {
	client       *http.Client
	BaseURL      *url.URL
	UserAgent    string
	Config       *Config
	BoardService *BoardService
}

func NewClient(config *Config, api_uri string) *Client {
	if config.APIVersion == "" {
		config.APIVersion = "v1"
	}
	baseURL, _ := url.Parse(api_uri)
	baseURL.Path = "api/" + config.APIVersion + "/"

	cl := http.DefaultClient

	c := &Client{
		client:    cl,
		BaseURL:   baseURL,
		UserAgent: userAgent,
		Config:    config,
	}

	c.initBoardService()
	return c
}

func (c *Client) initBoardService() error {
	bs := NewBoardService(c)
	if bs == nil {
		return errors.New("failed to init board service")
	}
	c.BoardService = bs
	return nil
}

func (c *Client) NewRequest(urlstr string, method string, body interface{}) (*http.Request, error) {
	//log.Printf("urlstr = %s\n", urlstr)
	rel, err := url.Parse(urlstr)
	if err != nil {
		return nil, err
	}

	//log.Printf("rel = %s\n", rel)
	//log.Printf("base_url = %s\n", c.BaseURL)
	u := c.BaseURL.ResolveReference(rel)
	//log.Printf("new req url = %s\n", u)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	//log.Printf("created new request with url %s\n", u.String())
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", mediaType)
	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

func (c *Client) Do(req *http.Request, into interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(into)
	return resp, err
}

func (c *Client) GetServerVersion() (*VersionInfo, error) {
	v := new(VersionInfo)
	req, err := c.NewRequest("version", "GET", nil)
	if err != nil {
		return nil, err
	}

	_, err = c.Do(req, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
