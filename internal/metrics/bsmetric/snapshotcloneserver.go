package bsmetric

import (
	"fmt"

	comm "github.com/opencurve/curve-manager/internal/metrics/common"
	"github.com/opencurve/curve-manager/internal/metrics/core"
)

const (
	SNAPSHOT_CLONE_STATUS           = "snapshotcloneserver_status"
	SNAPSHOT_CLONE_CONF_LISTEN_ADDR = "snapshotcloneserver_config_server_address"
	SNAPSHOT_CLONE_LEADER           = "active"
)

type ServiceStatus comm.ServiceStatus

func GetSnapShotCloneServerStatus() ([]ServiceStatus, string) {
	ret := []ServiceStatus{}
	size := len(core.GMetricClient.SnapShotCloneServerDummyAddr)
	if size == 0 {
		return ret, "no snapshotclone service address found"
	}
	results := make(chan comm.MetricResult, size)
	names := fmt.Sprintf("%s,%s,%s", comm.CURVEBS_VERSION, SNAPSHOT_CLONE_STATUS, SNAPSHOT_CLONE_CONF_LISTEN_ADDR)
	comm.GetBvarMetric(core.GMetricClient.SnapShotCloneServerDummyAddr, names, &results)

	count := 0
	var errors string
	for res := range results {
		if res.Err == nil {
			addr := ""
			v, e := comm.ParseBvarMetric(res.Result.(string))
			if e != nil {
				errors = fmt.Sprintf("%s; %s:%s", errors, res.Key, e.Error())
			} else {
				addr = comm.GetBvarConfMetricValue((*v)[SNAPSHOT_CLONE_CONF_LISTEN_ADDR])
			}
			ret = append(ret, ServiceStatus{
				Address: addr,
				Version: (*v)[comm.CURVEBS_VERSION],
				Online:  true,
				Leader:  (*v)[SNAPSHOT_CLONE_STATUS] == SNAPSHOT_CLONE_LEADER,
			})
		} else {
			errors = fmt.Sprintf("%s; %s:%s", errors, res.Key, res.Err.Error())
			ret = append(ret, ServiceStatus{
				Address: res.Key.(string),
				Version: "",
				Leader:  false,
				Online:  false,
			})
		}
		count += 1
		if count >= size {
			break
		}
	}
	return ret, errors
}
