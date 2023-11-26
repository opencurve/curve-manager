package baseHttp

import (
	"fmt"
	"net/url"
	"time"
)

type HttpResult struct {
	Key    interface{}
	Err    error
	Result interface{}
}

type BaseHttp struct {
	client     *resty.Client
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
		go func(addr string) {
			url := (&url.URL{
				Scheme: "http",
				Host:   addr,
				Path:   path,
			}).String()
			resp, err := cli.client.R().
				SetHeader("Connection", "Keep-Alive").
				SetHeader("Content-Type", "application/json").
				SetHeader("User-Agent", "curl/7.52.1").
				Execute("GET", url)
			results <- HttpResult{
				Key:    addr,
				Err:    err,
				Result: resp,
			}
		}(host)
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
