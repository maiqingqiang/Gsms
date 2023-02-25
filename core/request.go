package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"strings"
	"time"
)

type RequestInterface interface {
	Request(method, url string, data interface{}, options ...Option) ([]byte, error)
	Get(url string, data interface{}) ([]byte, error)
	Post(url string, data interface{}) ([]byte, error)
	GetWithUnmarshal(url string, data interface{}, v interface{}) (string, error)
	PostWithUnmarshal(url string, data interface{}, v interface{}) (string, error)
}

var _ RequestInterface = (*Request)(nil)

type Request struct {
	Timeout time.Duration
	BaseUrl string
}

func NewRequest(baseUrl string) *Request {
	return &Request{
		Timeout: time.Second * 5,
		BaseUrl: baseUrl,
	}
}

// Option http request option.
type Option func(req *http.Request)

// Request http request.
func (r *Request) Request(method, url string, data interface{}, options ...Option) ([]byte, error) {

	if !(strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https")) {
		url = r.BaseUrl + url
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if method == http.MethodGet {
		if data != nil {
			switch data.(type) {
			case map[string]string:
				q := req.URL.Query()
				for k, v := range data.(map[string]string) {
					q.Add(k, v)
				}
				req.URL.RawQuery = q.Encode()
			case string:
				req.URL.RawQuery = data.(string)
			case netUrl.Values:
				req.URL.RawQuery = data.(netUrl.Values).Encode()
			case []string:
				req.URL.RawQuery = strings.Join(data.([]string), "&")
			default:
				return nil, ErrRequestDataTypeError
			}
		}
	}

	client := &http.Client{
		Timeout: r.Timeout,
	}

	for _, option := range options {
		option(req)
	}

	if method == http.MethodPost {
		if data != nil {
			if req.Header.Get("Content-Type") == "application/json" {
				var body []byte

				switch data.(type) {
				case []byte:
					body = data.([]byte)
				case string:
					body = []byte(data.(string))
				default:
					body, err = json.Marshal(data)
					if err != nil {
						return nil, err
					}
				}

				req.Body = ioutil.NopCloser(bytes.NewReader(body))
			} else {
				var body string
				switch data.(type) {
				case map[string]string:
					q := netUrl.Values{}
					for k, v := range data.(map[string]string) {
						q.Add(k, v)
					}
					body = q.Encode()
				case string:
					body = data.(string)
				case netUrl.Values:
					body = data.(netUrl.Values).Encode()
				default:
					return nil, ErrRequestDataTypeError
				}

				req.Body = ioutil.NopCloser(strings.NewReader(body))
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Get request
func (r *Request) Get(url string, data interface{}) ([]byte, error) {
	return r.Request(http.MethodGet, url, data)
}

// GetWithUnmarshal get request and unmarshal response.
func (r *Request) GetWithUnmarshal(url string, data interface{}, v interface{}) (string, error) {
	body, err := r.Get(url, data)
	err = json.Unmarshal(body, v)
	if err != nil {
		return "", errors.New(fmt.Sprintf("json unmarshal error: %s body: %s", err.Error(), string(body)))
	}

	return string(body), nil
}

// Post request
func (r *Request) Post(url string, data interface{}) ([]byte, error) {
	return r.Request(http.MethodPost, url, data)
}

// PostWithUnmarshal post request and unmarshal response.
func (r *Request) PostWithUnmarshal(url string, data interface{}, v interface{}) (string, error) {
	body, err := r.Post(url, data)
	err = json.Unmarshal(body, v)
	if err != nil {
		return "", errors.New(fmt.Sprintf("json unmarshal error: %s body: %s", err.Error(), string(body)))
	}

	return string(body), nil
}
