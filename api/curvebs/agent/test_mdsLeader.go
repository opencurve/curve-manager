package agent

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/opencurve/curve-manager/internal/common"
	"net/url"
	"testing"
)

func TestMdsClient_ListPhysicalPool_http(t *testing.T) {

}
func TestGetCurrentClusterServicesAddr() (clusterServicesAddr, error) {
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
