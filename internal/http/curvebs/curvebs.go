package curvebs

import (
	"github.com/go-resty/resty/v2"
	"github.com/opencurve/curve-manager/internal/common"
	"net/url"
	"strings"
)

var (
	GMdsClient *MdsClient
)

const (
	CURVEBS_MDS_ADDRESS = "mds.address"

	DEFAULT_RPC_TIMEOUT_MS  = 500
	DEFAULT_RPC_RETRY_TIMES = 3
)

func Init(cfg map[string]string) {
	addrs := findLeader(cfg)
	GMdsClient = NewMdsClient(MdsClientOption{
		TimeoutMs:  DEFAULT_RPC_TIMEOUT_MS,
		RetryTimes: DEFAULT_RPC_RETRY_TIMES,
		Addrs:      strings.Split(addrs, common.CURVEBS_ADDRESS_DELIMITER),
	})
}

func findLeader(cfg map[string]string) string {
	mds_addr := cfg[CURVEBS_MDS_ADDRESS]
	Addrs := strings.Split(mds_addr, common.CURVEBS_ADDRESS_DELIMITER)
	for _, addr := range Addrs {
		httpClient := common.GetHttpClient()
		url := (&url.URL{
			Scheme: "http",
			Host:   addr,
			Path:   "/",
		}).String()
		resp, err := resty.NewWithClient(httpClient).R().
			SetHeader("Connection", "Keep-Alive").
			SetHeader("Content-Type", "application/json").
			SetHeader("User-Agent", "Curve-Manager").
			Execute("GET", url)
		if err != nil {
			continue
		}
		if resp.Body() != nil {
			return addr
		}

	}
	return mds_addr
}
