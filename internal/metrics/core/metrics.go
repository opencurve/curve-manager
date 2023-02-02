package core

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/pigeon"
)

const (
	CURVEBS_MONITOR_ADDRESS              = "monitor.prometheus.address"
	CURVEBS_ETCD_ADDRESS                 = "etcd.address"
	CURVEBS_MDS_DUMMY_ADDRESS            = "mds.dummy.address"
	CURVEBS_SNAPSHOT_CLONE_DUMMY_ADDRESS = "snapshot.clone.dummy.address"
)

type metricClient struct {
	client                       *resty.Client
	PromeAddr                    string
	EtcdAddr                     []string
	MdsDummyAddr                 []string
	SnapShotCloneServerDummyAddr []string
}

var (
	GMetricClient *metricClient
)

func Init(cfg *pigeon.Configure) error {
	GMetricClient = &metricClient{}
	// init http client
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

	GMetricClient.client = resty.NewWithClient(httpClient)
	addr := cfg.GetConfig().GetString(CURVEBS_MONITOR_ADDRESS)
	if len(addr) == 0 {
		return fmt.Errorf("no cluster monitor address found")
	}
	GMetricClient.PromeAddr = addr

	addr = cfg.GetConfig().GetString(CURVEBS_ETCD_ADDRESS)
	if len(addr) == 0 {
		return fmt.Errorf("no cluster etcd address found")
	}
	GMetricClient.EtcdAddr = strings.Split(addr, common.CURVEBS_ADDRESS_DELIMITER)

	addr = cfg.GetConfig().GetString(CURVEBS_MDS_DUMMY_ADDRESS)
	if len(addr) == 0 {
		return fmt.Errorf("no cluster mds dummy address found")
	}
	GMetricClient.MdsDummyAddr = strings.Split(addr, common.CURVEBS_ADDRESS_DELIMITER)

	addr = cfg.GetConfig().GetString(CURVEBS_SNAPSHOT_CLONE_DUMMY_ADDRESS)
	if len(addr) != 0 {
		GMetricClient.SnapShotCloneServerDummyAddr = strings.Split(addr, common.CURVEBS_ADDRESS_DELIMITER)
	}

	return nil
}

func (cli *metricClient) GetMetricFromService(host, path string) (interface{}, error) {
	url := (&url.URL{
		Scheme: "http",
		Host:   host,
		Path:   path,
	}).String()

	resp, err := cli.client.R().
		SetHeader("Connection", "Keep-Alive").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "curl/7.52.1").
		Execute("GET", url)
	if err != nil {
		return nil, fmt.Errorf("get bvar metric failed, addr: %s, err: %v", host, err)
	} else if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get bvar metric failed, addr: %s, status = %s", host, resp.Status())
	}
	return resp.String(), nil
}

func (cli *metricClient) GetMetricFromPrometheus(path, queryKey, queryValue string, ret interface{}) error {
	url := (&url.URL{
		Scheme:   "http",
		Host:     cli.PromeAddr,
		Path:     path,
		RawQuery: fmt.Sprintf("%s=%s", queryKey, queryValue),
	}).String()

	resp, err := cli.client.R().
		SetHeader("Connection", "Keep-Alive").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Curve-Manager").
		SetResult(ret).
		Execute("GET", url)
	if err != nil {
		return fmt.Errorf("get prometheus metric failed: %v", err)
	} else if resp.StatusCode() != 200 {
		return fmt.Errorf("get prometheus metric failed, status = %s",
			resp.Status())
	}
	return nil
}
