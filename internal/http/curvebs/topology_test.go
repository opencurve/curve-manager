package curvebs

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/opencurve/curve-manager/internal/common"
	"net/url"
	"strings"
	"testing"
)

var (
	clientOption MdsClientOption = MdsClientOption{
		TimeoutMs:  500,
		RetryTimes: 3,
		Addrs:      []string{"192.168.170.138:6702"},
	}
)

const (
	CLUSTER_SERVICES_ADDRESS = "cluster.service.addr"
)

type admHttpResponse struct {
	ErrorCode string              `json:"errorCode"`
	ErrorMsg  string              `json:"errorMsg"`
	Data      clusterServicesAddr `json:"data"`
}
type clusterServicesAddr struct {
	ClusterId int               `json:"clusterId"`
	Addrs     map[string]string `json:"addrs"`
}

func TestMdsClient_ListPhysicalPool_http(t *testing.T) {
	info, err := GetCurrentClusterServicesAddr()
	if err != nil {

	}
	mds_addr := info.Addrs[CURVEBS_MDS_ADDRESS]
	if mds_addr != "" {
		Addrs := strings.Split(mds_addr, common.CURVEBS_ADDRESS_DELIMITER)
		for _, addr := range Addrs {
			httpClient := common.GetHttpClient()
			url := (&url.URL{
				Scheme:   "http",
				Host:     addr,
				Path:     "/",
				RawQuery: fmt.Sprintf("%s=%s", "method", CLUSTER_SERVICES_ADDRESS),
			}).String()

			resp, err := resty.NewWithClient(httpClient).R().
				SetHeader("Connection", "Keep-Alive").
				SetHeader("Content-Type", "application/json").
				SetHeader("User-Agent", "Curve-Manager").
				Execute("GET", url)
			if resp.Body() != nil {

			}
			if err != nil {

			}
		}
	}

}

func GetCurrentClusterServicesAddr() (clusterServicesAddr, error) {
	ret := clusterServicesAddr{}
	httpClient := common.GetHttpClient()
	url := (&url.URL{
		Scheme:   "http",
		Host:     "127.0.0.1:11000",
		Path:     "/",
		RawQuery: fmt.Sprintf("%s=%s", "method", CLUSTER_SERVICES_ADDRESS),
	}).String()

	resp, err := resty.NewWithClient(httpClient).R().
		SetHeader("Connection", "Keep-Alive").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Curve-Manager").
		Execute("GET", url)
	if err != nil {
		return ret, fmt.Errorf("getClusterServicesAddr failed: %v", err)
	}

	respStruct := admHttpResponse{}
	err = json.Unmarshal([]byte(resp.String()), &respStruct)
	if err != nil {
		return ret, fmt.Errorf("Unmarshal getClusterServicesAddr response failed, resp = %s, err = %v", resp.String(), err)
	}
	return respStruct.Data, nil
}
