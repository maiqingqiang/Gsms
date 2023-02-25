package core

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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
	resp, err := r.do(method, url, data, options...)
	if err != nil {
		return nil, err
	}

	return r.readBody(resp)
}

func (r *Request) readBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (r *Request) do(method string, url string, data interface{}, options ...Option) (*http.Response, error) {
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
						return nil, ErrRequestDataTypeError
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
	return resp, nil
}

// Get request
func (r *Request) Get(url string, data interface{}) ([]byte, error) {
	return r.Request(http.MethodGet, url, data)
}

// GetWithUnmarshal get request and unmarshal response.
func (r *Request) GetWithUnmarshal(url string, data interface{}, v interface{}) (string, error) {
	resp, err := r.do(http.MethodGet, url, data)

	body, err := r.readBody(resp)
	if err != nil {
		return "", err
	}

	bodyStr := string(body)

	if strings.Contains(resp.Header.Get("Content-Type"), "json") || strings.Contains(resp.Header.Get("Content-Type"), "javascript") {
		err = json.Unmarshal(body, v)
		if err != nil {
			return "", errors.New(fmt.Sprintf("json unmarshal error: %s, body: %s", err, bodyStr))
		}
	} else if strings.Contains(resp.Header.Get("Content-Type"), "xml") {
		err = xml.Unmarshal(body, v)
		if err != nil {
			return "", errors.New(fmt.Sprintf("xml unmarshal error: %s, body: %s", err, bodyStr))
		}
	} else {
		if resp.StatusCode != http.StatusOK {
			return "", errors.New(fmt.Sprintf("not support unmarshal content type: %s, status code: %d, body: %s", resp.Header.Get("Content-Type"), resp.StatusCode, bodyStr))
		}
	}

	return bodyStr, nil
}

// Post request
func (r *Request) Post(url string, data interface{}) ([]byte, error) {
	return r.Request(http.MethodPost, url, data)
}

// PostWithUnmarshal post request and unmarshal response.
func (r *Request) PostWithUnmarshal(url string, data interface{}, v interface{}) (string, error) {
	resp, err := r.do(http.MethodPost, url, data)

	body, err := r.readBody(resp)
	if err != nil {
		return "", err
	}

	err = r.unmarshal(resp, body, v)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// unmarshal response.
func (r *Request) unmarshal(resp *http.Response, body []byte, v interface{}) error {
	bodyStr := string(body)

	if strings.Contains(resp.Header.Get("Content-Type"), "json") || strings.Contains(resp.Header.Get("Content-Type"), "javascript") {
		err := json.Unmarshal(body, v)
		if err != nil {
			return errors.New(fmt.Sprintf("json unmarshal error: %s, body: %s", err, bodyStr))
		}
	} else if strings.Contains(resp.Header.Get("Content-Type"), "xml") {
		err := xml.Unmarshal(body, v)
		if err != nil {
			return errors.New(fmt.Sprintf("xml unmarshal error: %s, body: %s", err, bodyStr))
		}
	} else {
		if resp.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("not support unmarshal content type: %s, status code: %d, body: %s", resp.Header.Get("Content-Type"), resp.StatusCode, bodyStr))
		}
	}
	return nil
}
