package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ClientInterface interface {
	GetWithUnmarshal(api string, data interface{}, v ResponseInterface) (string, error)
	PostFormWithUnmarshal(api string, data string, v ResponseInterface) (string, error)
}

type ResponseInterface interface {
	Unmarshal(data []byte) error
}

type Client struct {
	HttpClient *http.Client
	Debug      bool
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{HttpClient: httpClient}
}

func (c *Client) Do(r *http.Request) (*http.Response, error) {
	if r.Header.Get("Content-Type") == "" {
		r.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	if c.Debug {
		log.Printf("Request: %s %s %s", r.Method, r.URL.String(), r.Body)
	}

	return c.HttpClient.Do(r)
}

func (c *Client) Post(api string, data map[string]interface{}) (*http.Response, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, api, bytes.NewReader(d))

	return c.Do(r)
}

func (c *Client) PostForm(api string, data string) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodPost, api, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.Do(r)
}

// PostFormWithUnmarshal is a helper function to post form data and unmarshal the response
func (c *Client) PostFormWithUnmarshal(api string, data string, v ResponseInterface) (string, error) {
	resp, err := c.PostForm(api, data)
	if err != nil {
		return "", err
	}

	body, err := c.Unmarshal(resp, v)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) Get(api string, data interface{}) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		return nil, err
	}

	switch data.(type) {
	case map[string]string:
		q := r.URL.Query()
		for k, v := range data.(map[string]string) {
			q.Add(k, v)
		}
		r.URL.RawQuery = q.Encode()
	case url.Values:
		r.URL.RawQuery = data.(url.Values).Encode()
	default:
		return nil, ErrRequestDataTypeError
	}

	return c.Do(r)
}

func (c *Client) GetWithUnmarshal(api string, data interface{}, v ResponseInterface) (string, error) {
	resp, err := c.Get(api, data)
	if err != nil {
		return "", err
	}

	body, err := c.Unmarshal(resp, v)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) Unmarshal(resp *http.Response, v ResponseInterface) ([]byte, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = v.Unmarshal(body)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if err != nil {
			return nil, fmt.Errorf("http status code: %d, error: %s", resp.StatusCode, body)
		}
		return nil, err
	}

	return body, nil
}
