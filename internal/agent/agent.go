package agent

import (
	"fmt"
	"sort"
	"time"

	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/metrics/bsmetric"
	bsrpc "github.com/opencurve/curve-manager/internal/rpc/curvebs"
)

const (
	ORDER_BY_ID     = "id"
	ORDER_BY_CTIME  = "ctime"
	ORDER_BY_LENGTH = "length"
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
		return nil, err
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

func GetClusterStatus() (interface{}, error) {
	healthy := true
	// 1. check service status
	// etcd, mds, snapshotcloneserver
	size := 3
	ret := make(chan bool, size)

	go func() {
		ret <- checkServiceHealthy(ETCD_SERVICE)
	}()

	go func() {
		ret <- checkServiceHealthy(MDS_SERVICE)
	}()

	go func() {
		ret <- checkServiceHealthy(SNAPSHOT_CLONE_SERVER_SERVICE)
	}()
	count := 0
	for res := range ret {
		healthy = healthy && res
		count += 1
		if count >= size {
			break
		}
	}
	// 2. check copyset status
	// 2.1 get chunkservers in cluster
	chunkservers, err := bsrpc.GMdsClient.GetChunkServerInCluster()
	if err != nil {
		healthy = false
	}

	// TODO
	return chunkservers, nil
}

func GetClusterSpace() (interface{}, error) {
	var result Space
	// get logical pools form mds
	pools, err := bsrpc.GMdsClient.ListLogicalPool()
	if err != nil {
		return nil, err
	}

	var poolInfos []PoolInfo
	for _, pool := range pools {
		var info PoolInfo
		info.Name = pool.Name
		poolInfos = append(poolInfos, info)
	}

	err = getPoolSpace(&poolInfos)
	if err != nil {
		return nil, err
	}
	for _, info := range poolInfos {
		result.Total += info.Space.Total
		result.Alloc += info.Space.Alloc
		result.CanRecycled += info.Space.CanRecycled
	}
	return &result, nil
}

func GetClusterPerformance() (interface{}, error) {
	return bsmetric.GetClusterPerformance()
}

func ListTopology() (interface{}, error) {
	var result []Pool
	logicalPools, err := bsrpc.GMdsClient.ListLogicalPool()
	if err != nil {
		return nil, err
	}
	for _, lp := range logicalPools {
		var pool Pool
		pool.Id = lp.Id
		pool.physicalPoolId = lp.PhysicalPoolId
		pool.Name = lp.Name
		pool.Type = lp.Type
		result = append(result, pool)
	}
	err = listPoolZone(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func ListLogicalPool() (interface{}, error) {
	var result []PoolInfo
	// get info from mds
	pools, err := bsrpc.GMdsClient.ListLogicalPool()
	if err != nil {
		return nil, err
	}

	for _, pool := range pools {
		var info PoolInfo
		info.Id = pool.Id
		info.Name = pool.Name
		info.PhysicalPoolId = pool.PhysicalPoolId
		info.Type = pool.Type
		info.CreateTime = pool.CreateTime
		info.AllocateStatus = pool.AllocateStatus
		info.ScanEnable = pool.ScanEnable
		result = append(result, info)
	}

	// get info from monitor
	err = getPoolItemNum(&result)
	if err != nil {
		return nil, err
	}

	err = getPoolSpace(&result)
	if err != nil {
		return nil, err
	}

	err = getPoolPerformance(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func sortFile(files []bsrpc.FileInfo, orderKey string) {
	sort.Slice(files, func(i, j int) bool {
		switch orderKey {
		case ORDER_BY_CTIME:
			itime, _ := time.Parse(common.TIME_FORMAT, files[i].Ctime)
			jtime, _ := time.Parse(common.TIME_FORMAT, files[j].Ctime)
			return itime.Unix() < jtime.Unix()
		case ORDER_BY_LENGTH:
			return files[i].Length < files[j].Length
		}
		return files[i].Id < files[j].Id
	})
}

func ListVolume(size, page uint32, path, key string) (interface{}, error) {
	// get root auth info
	authInfo, err := bsmetric.GetAuthInfoOfRoot()
	if err != "" {
		return nil, fmt.Errorf(err)
	}

	// create signature
	date := time.Now().UnixMicro()
	str2sig := common.GetString2Signature(date, authInfo.UserName)
	sig := common.CalcString2Signature(str2sig, authInfo.PassWord)

	fileInfos, e := bsrpc.GMdsClient.ListDir(path, authInfo.UserName, sig, uint64(date))
	if e != nil {
		return nil, e
	}

	if len(fileInfos) == 0 {
		return nil, nil
	}

	sortFile(fileInfos, key)
	length := uint32(len(fileInfos))
	start := (page - 1) * size
	var end uint32
	if page*size > length {
		end = length
	} else {
		end = page * size
	}
	return fileInfos[start:end], nil
}
