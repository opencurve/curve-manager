package metrics

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/opencurve/pigeon"
)

const (
	CURVEBS_MONITOR_ADDRESS = "monitor.prometheus.addr"
)

var (
	client *resty.Client
	addr   string
)

func init() {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxConnsPerHost:       100,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}

	client = resty.NewWithClient(httpClient)
}

func Init(cfg *pigeon.Configure) error {
	addr = cfg.GetConfig().GetString(CURVEBS_MONITOR_ADDRESS)
	if len(addr) == 0 {
		return fmt.Errorf("no cluster monitor address found")
	}
	return nil
}

func sendHttp(method, url string, ret interface{}) error {
	resp, err := client.R().
		SetHeader("Connection", "Keep-Alive").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Curve-Manager").
		SetResult(ret).
		Execute(method, url)
	if err != nil {
		return fmt.Errorf("get metric failed: %v", err)
	} else if resp.StatusCode() != 200 {
		return fmt.Errorf("get metric failed, status = %s",
			resp.Status())
	}
	return nil
}

func getMetrics(path, queryKey, queryValue string, ret interface{}) error {
	return sendHttp("GET",
		(&url.URL{
			Scheme:   "http",
			Host:     addr,
			Path:     path,
			RawQuery: fmt.Sprintf("%s=%s", queryKey, queryValue),
		}).String(), ret)
}

func parseVectorMetric(info *QueryResponseOfVector, key string, isValue bool) *map[string]string {
	ret := make(map[string]string)
	if info == nil {
		return &ret
	}
	for _, metric := range info.Data.Result {
		if !isValue {
			ret[metric.Metric["instance"]] = metric.Metric[key]
		} else {
			ret[metric.Metric["instance"]] = metric.Value[1].(string)
		}
	}
	return &ret
}

func GetEtcdStatus() (*EtcdStatus, error) {
	var verResult, leaderResult QueryResponseOfVector
	var verErr, leaderErr error
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		verErr = getMetrics(VECTOR_METRIC_PATH, VECTOR_METRIC_QUERY_KEY, ETCD_CLUSTER_VERSION_NAME, &verResult)
	}()

	go func() {
		defer wg.Done()
		leaderErr = getMetrics(VECTOR_METRIC_PATH, VECTOR_METRIC_QUERY_KEY, ETCD_SERVER_IS_LEADER_NAME, &leaderResult)
	}()

	wg.Wait()
	var ret EtcdStatus
	if verErr != nil {
		return nil, verErr
	}
	// get etcd version
	versions := parseVectorMetric(&verResult, ETCD_CLUSTER_VERSION, false)
	if len(*versions) == 0 {
		return nil, fmt.Errorf("no etcd cluster version found")
	}
	for _, v := range *versions {
		ret.Version = v
		break
	}

	if leaderErr != nil {
		return nil, leaderErr
	}
	// get etcd leader
	leaders := parseVectorMetric(&leaderResult, "", true)
	if len(*leaders) == 0 {
		return nil, fmt.Errorf("no etcd role found")
	}
	for k, v := range *leaders {
		if v == "1" {
			ret.Leader = append(ret.Leader, k)
		} else {
			ret.Follower = append(ret.Follower, k)
		}
	}
	return &ret, nil
}
