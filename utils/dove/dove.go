package dove

import (
	"encoding/json"
	"fmt"
	"github.com/maiqingqiang/gsms"
	"io"
	"net/http"
)

type StatusCodeJudger func(statusCode int) error
type Unmarshal func(data []byte, v interface{}) error

type Dove struct {
	client           *http.Client
	statusCodeJudger StatusCodeJudger
	unmarshal        Unmarshal
	logger           gsms.Logger
}

func New(opts ...Option) *Dove {
	dove := &Dove{
		client: http.DefaultClient,
		statusCodeJudger: func(statusCode int) error {
			if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
				return fmt.Errorf("request status code: %d", statusCode)
			}

			return nil
		},
		unmarshal: json.Unmarshal,
	}

	for _, opt := range opts {
		opt(dove)
	}

	return dove
}

func (d *Dove) Request(method, url string, header http.Header, data io.Reader, response interface{}) error {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return err
	}

	if header != nil {
		req.Header = header
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	d.logger.Infof("request url: %s , method: %s , header: %+v , reqBody: %s", req.URL, req.Method, req.Header, req.Body)

	resp, err := d.client.Do(req)
	if err != nil {
		d.logger.Warnf("request failed: %v", err)
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		d.logger.Warnf("io.ReadAll failed: %v", err)
		return err
	}

	d.logger.Infof("response status code: %d , body: %s", resp.StatusCode, body)

	err = d.statusCodeJudger(resp.StatusCode)
	if err != nil {
		return err
	}

	err = d.unmarshal(body, response)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dove) Get(url string, data io.Reader, response interface{}) error {
	return d.Request(http.MethodGet, url, nil, data, response)
}

func (d *Dove) Post(url string, data io.Reader, response interface{}) error {
	return d.Request(http.MethodPost, url, nil, data, response)
}

func (d *Dove) PostForm(url string, data io.Reader, response interface{}) error {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")

	return d.Request(http.MethodPost, url, header, data, response)
}
