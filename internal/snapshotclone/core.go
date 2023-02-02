package snapshotclone

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

type snapshotCloneClient struct {
	client                 *resty.Client
	SnapShotCloneProxyAddr []string
}

var (
	GSnapshotCloneClient *snapshotCloneClient
)

const (
	CURVEBS_SNAPSHOT_CLONE_PROXY_ADDRESS = "snapshot.clone.proxy.address"
)

func Init(cfg *pigeon.Configure) error {
	GSnapshotCloneClient = &snapshotCloneClient{}
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

	GSnapshotCloneClient.client = resty.NewWithClient(httpClient)
	addr := cfg.GetConfig().GetString(CURVEBS_SNAPSHOT_CLONE_PROXY_ADDRESS)
	if len(addr) != 0 {
		GSnapshotCloneClient.SnapShotCloneProxyAddr = strings.Split(addr, common.CURVEBS_ADDRESS_DELIMITER)
	}
	if len(GSnapshotCloneClient.SnapShotCloneProxyAddr) == 0 {
		return fmt.Errorf("have no valide snapshotclone proxy addr")
	}

	return nil
}

func (cli *snapshotCloneClient) sendHttp2SnapshotClone(queryParam string) (string, error) {
	url := (&url.URL{
		Scheme:   "http",
		Host:     cli.SnapShotCloneProxyAddr[0],
		Path:     SERVICE_ROUTER,
		RawQuery: queryParam,
	}).String()

	resp, err := cli.client.R().
		SetHeader("Connection", "Keep-Alive").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Curve-Manager").
		Execute("GET", url)
	if err != nil {
		return "", fmt.Errorf("get snapshot failed: %v", err)
	} else if resp.StatusCode() != 200 {
		return "", fmt.Errorf("get snapshot failed, status = %s",
			resp.Status())
	}
	return resp.String(), nil
}
