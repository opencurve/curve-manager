package bsmetric

import (
	"fmt"

	"github.com/opencurve/curve-manager/internal/metrics/core"
)

const (
	MDS_STATUS           = "mds_status"
	MDS_LEADER           = "leader"
	MDS_FOLLOWER         = "follower"
	MDS_CONF_LISTEN_ADDR = "mds_config_mds_listen_addr"
)

func GetMdsStatus() (*[]MdsStatus, string) {
	size := len(core.GMetricClient.MdsDummyAddr)
	results := make(chan metricResult, size)
	names := fmt.Sprintf("%s,%s,%s", CURVEBS_VERSION, MDS_STATUS, MDS_CONF_LISTEN_ADDR)
	getBvarMetric(core.GMetricClient.MdsDummyAddr, names, &results)

	count := 0
	var errors string
	var ret []MdsStatus
	for res := range results {
		if res.err == nil {
			addr := ""
			v, e := parseBvarMetric(res.result.(string))
			if e != nil {
				errors = fmt.Sprintf("%s; %s:%s", errors, res.key, e.Error())
			} else {
				addr = getBvarConfMetricValue((*v)[MDS_CONF_LISTEN_ADDR])
			}
			ret = append(ret, MdsStatus{
				Address: addr,
				Version: (*v)[CURVEBS_VERSION],
				Online:  true,
				Leader:  (*v)[MDS_STATUS] == MDS_LEADER,
			})
		} else {
			errors = fmt.Sprintf("%s; %s:%s", errors, res.key, res.err.Error())
			ret = append(ret, MdsStatus{
				Address: res.key.(string),
				Version: "",
				Leader:  false,
				Online:  false,
			})
		}
		count = count + 1
		if count >= size {
			break
		}
	}
	return &ret, errors
}
