package agent

import (
	"fmt"

	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/metrics/bsmetric"
	bsrpc "github.com/opencurve/curve-manager/internal/rpc/curvebs"
)

func GetEtcdStatus() (interface{}, string) {
	return bsmetric.GetEtcdStatus()
}

func GetMdsStatus() (interface{}, string) {
	return bsmetric.GetMdsStatus()
}

func GetSnapShotCloneServerStatus() (interface{}, string) {
	return bsmetric.GetSnapShotCloneServerStatus()
}

func GetChunkServerStatus() (interface{}, error) {
	var result ChunkServerStatus
	// get chunkserver form mds
	chunkservers, err := bsrpc.GMdsClient.GetChunkServerInCluster()
	if err != nil {
		return nil, fmt.Errorf("GetChunkServerInCluster failed, %s", err)
	}

	online := 0
	var endponits []string
	for _, cs := range chunkservers {
		if cs.OnlineStatus == bsrpc.ONLINE_STATUS {
			online += 1
		}
		endponits = append(endponits, fmt.Sprintf("%s:%d", cs.HostIp, cs.Port))
	}
	result.TotalNum = len(chunkservers)
	result.OnlineNum = online

	// get version form metric
	versions, err := bsmetric.GetChunkServerVersion(&endponits)
	if err != nil {
		return nil, err
	}
	for k, v := range *versions {
		result.Versions = append(result.Versions, VersionNum{
			Version: k,
			Number:  v,
		})
	}
	return &result, nil
}

func checkServiceHealthy(name string) common.QueryResult {
	var ret common.QueryResult
	var status []bsmetric.ServiceStatus
	var err string
	ret.Key = name
	switch name {
	case ETCD_SERVICE:
		status, err = bsmetric.GetEtcdStatus()
	case MDS_SERVICE:
		status, err = bsmetric.GetMdsStatus()
	case SNAPSHOT_CLONE_SERVER_SERVICE:
		status, err = bsmetric.GetSnapShotCloneServerStatus()
	default:
		ret.Result = false
		ret.Err = fmt.Errorf("Invalid service name")
		return ret
	}

	if err != "" {
		ret.Err = fmt.Errorf(err)
		ret.Result = false
		return ret
	}

	leaderNum := 0
	hasOffline := false
	for _, s := range status {
		if s.Leader {
			leaderNum += 1
		}
		if !s.Online {
			hasOffline = true
		}
	}

	ret.Result = leaderNum == 1 && !hasOffline
	return ret
}