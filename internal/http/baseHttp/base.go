package baseHttp

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
	"time"
)

type HttpResult struct {
	Key    interface{}
	Err    error
	Result interface{}
}

type BaseHttp struct {
	Client     *resty.Client
	Timeout    time.Duration
	RetryTimes uint32
}

var (
	GMetricClient *BaseHttp
)

func (cli *BaseHttp) SendHTTP(host []string, path string) *HttpResult {

	size := len(host)
	if size == 0 {
		return &HttpResult{
			Key:    "",
			Err:    fmt.Errorf("empty addr"),
			Result: nil,
		}
	}
	results := make(chan HttpResult, size)
	for _, host := range host {
		url := (&url.URL{
			Scheme: "http",
			Host:   host,
			Path:   path,
		}).String()

		resp, err := cli.Client.R().
			SetHeader("Connection", "Keep-Alive").
			SetHeader("Content-Type", "application/json").
			SetHeader("User-Agent", "curl/7.52.1").
			Execute("GET", url)
		results <- HttpResult{
			Key:    host,
			Err:    err,
			Result: resp,
		}
	}
	var count = 0
	var httpErr string
	for res := range results {
		if res.Err == nil {
			return &res
		}
		count++
		httpErr = fmt.Sprintf("%s;%s:%s", httpErr, res.Key, res.Err.Error())
		if count >= size {
			break
		}

	}
	return &HttpResult{
		Key:    "",
		Err:    fmt.Errorf(httpErr),
		Result: nil,
	}

}

func (cli *BaseHttp) SendHTTPByPost(host []string, path string, body any) *HttpResult {

	size := len(host)
	if size == 0 {
		return &HttpResult{
			Key:    "",
			Err:    fmt.Errorf("empty addr"),
			Result: nil,
		}
	}
	results := make(chan HttpResult, size)
	for _, host := range host {
		url := (&url.URL{
			Scheme: "http",
			Host:   host,
			Path:   path,
		}).String()

		resp, err := cli.Client.R().
			SetHeader("Connection", "Keep-Alive").
			SetHeader("Content-Type", "application/json").
			SetHeader("User-Agent", "curl/7.52.1").
			SetBody(body).
			Execute("Post", url)
		results <- HttpResult{
			Key:    host,
			Err:    err,
			Result: resp,
		}
	}
	var count = 0
	var httpErr string
	for res := range results {
		if res.Err == nil {
			return &res
		}
		count++
		httpErr = fmt.Sprintf("%s;%s:%s", httpErr, res.Key, res.Err.Error())
		if count >= size {
			break
		}

	}
	return &HttpResult{
		Key:    "",
		Err:    fmt.Errorf(httpErr),
		Result: nil,
	}

}
