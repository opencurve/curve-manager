package bsmetric

import (
	"fmt"

	"github.com/opencurve/curve-manager/internal/metrics/core"
)

func GetMdsStatus() (*[]MdsStatus, error) {
	size := len(core.GMetricClient.MdsDummyAddr)
	results := make(chan metricResult, size)
	names := fmt.Sprintf("%s,%s",CURVEBS_VERSION, MDS_STATUS)
	getBvarMetric(core.GMetricClient.MdsDummyAddr, names, &results)

	count := 0
	var errors string
	var ret []MdsStatus
	for res := range results {
		if res.err == nil {
			v, e := parseBvarMetric(res.result.(string))
			if e != nil {
				errors = fmt.Sprintf("%s; %s:%s", errors, res.addr, e.Error())
			}
			ret = append(ret, MdsStatus{
				Address: res.addr,
				Version: (*v)[CURVEBS_VERSION],
				Online: true,
				Leader: (*v)[MDS_STATUS] == MDS_LEADER,
			})
		} else {
			errors = fmt.Sprintf("%s; %s:%s", errors, res.addr, res.err.Error())
			ret = append(ret, MdsStatus{
				Address: res.addr,
				Version: "",
				Leader: false,
				Online: false,
			})
		}
		count = count + 1
		if count >= size {
			break
		}
	}
	return &ret, fmt.Errorf(errors)
}
