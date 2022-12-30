package bsmetric

import (
	"fmt"

	"github.com/opencurve/curve-manager/internal/metrics/core"
)

const (
	SNAPSHOT_CLONE_STATUS           = "snapshotcloneserver_status"
	SNAPSHOT_CLONE_CONF_LISTEN_ADDR = "snapshotcloneserver_config_server_address"
	SNAPSHOT_CLONE_LEADER           = "active"
)

func GetSnapShotCloneServerStatus() (*[]SnapShotCloneServerStatus, string) {
	size := len(core.GMetricClient.SnapShotCloneServerDummyAddr)
	results := make(chan metricResult, size)
	names := fmt.Sprintf("%s,%s,%s", CURVEBS_VERSION, SNAPSHOT_CLONE_STATUS, SNAPSHOT_CLONE_CONF_LISTEN_ADDR)
	getBvarMetric(core.GMetricClient.SnapShotCloneServerDummyAddr, names, &results)

	count := 0
	var errors string
	var ret []SnapShotCloneServerStatus
	for res := range results {
		if res.err == nil {
			addr := ""
			v, e := parseBvarMetric(res.result.(string))
			if e != nil {
				errors = fmt.Sprintf("%s; %s:%s", errors, res.key, e.Error())
			} else {
				addr = getBvarConfMetricValue((*v)[SNAPSHOT_CLONE_CONF_LISTEN_ADDR])
			}
			ret = append(ret, SnapShotCloneServerStatus{
				Address: addr,
				Version: (*v)[CURVEBS_VERSION],
				Online:  true,
				Leader:  (*v)[SNAPSHOT_CLONE_STATUS] == SNAPSHOT_CLONE_LEADER,
			})
		} else {
			errors = fmt.Sprintf("%s; %s:%s", errors, res.key, res.err.Error())
			ret = append(ret, SnapShotCloneServerStatus{
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
