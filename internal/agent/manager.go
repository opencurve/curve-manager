package agent

import (
	"github.com/opencurve/curve-manager/internal/metrics/bsmetric"
	"github.com/opencurve/curve-manager/internal/rpc/baserpc"
	"github.com/opencurve/curve-manager/internal/rpc/curvebs/common"
	"github.com/opencurve/curve-manager/internal/rpc/curvebs/mds"
)

func GetEtcdStatus() (interface{}, error) {
	return bsmetric.GetEtcdStatus()
}

func GetMdsStatus() (interface{}, string) {
	return bsmetric.GetMdsStatus()
}

func GetSnapShotCloneServerStatus() (interface{}, string) {
	return bsmetric.GetSnapShotCloneServerStatus()
}

func GetChunkServerStatus() (interface{}, string) {
	return nil, ""
}

// get chunkservers of server concurrently
func listChunkServer(pools *[]Pool, size int) error {
	ret := make(chan baserpc.RpcResult, size)
	count := 0
	for pIndex, pool := range *pools {
		for zIndex, zone := range pool.Zones {
			for sIndex, server := range zone.Servers {
				go func(id uint32, addr *Server) {
					chunkservers, err := mds.GMdsClient.ListChunkServer(id)
					ret <- baserpc.RpcResult{
						Key:    addr,
						Result: chunkservers,
						Err:    err,
					}
				}(server.Id, &(*pools)[pIndex].Zones[zIndex].Servers[sIndex])
			}
		}
	}
	for res := range ret {
		if res.Err != nil {
			return res.Err
		}
		for _, cs := range res.Result.([]common.ChunkServer) {
			res.Key.(*Server).ChunkServers = append(res.Key.(*Server).ChunkServers, cs)
		}
		count = count + 1
		if count >= size {
			break
		}
	}
	return nil
}

// get servers of zones concurrently
func listZoneServer(pools *[]Pool, size int) error {
	ret := make(chan baserpc.RpcResult, size)
	count := 0
	number := 0
	for pIndex, pool := range *pools {
		for zIndex, zone := range pool.Zones {
			go func(id uint32, addr *Zone) {
				servers, err := mds.GMdsClient.ListZoneServer(id)
				ret <- baserpc.RpcResult{
					Key:    addr,
					Result: servers,
					Err:    err,
				}
			}(zone.Id, &(*pools)[pIndex].Zones[zIndex])
		}
	}
	for res := range ret {
		if res.Err != nil {
			return res.Err
		}
		for _, s := range res.Result.([]common.Server) {
			var server Server
			server.Id = s.Id
			server.Hostname = s.HostName
			server.InternalIp = s.InternalIp
			server.InternalPort = s.InternalPort
			server.ExternalIp = s.ExternalIp
			server.ExternalPort = s.ExternalPort
			res.Key.(*Zone).Servers = append(res.Key.(*Zone).Servers, server)
			number = number + 1
		}
		count = count + 1
		if count >= size {
			break
		}
	}
	return listChunkServer(pools, number)
}

// get zones of pools concurrently
func listPoolZone(pools *[]Pool) error {
	size := len(*pools)
	ret := make(chan baserpc.RpcResult, size)
	count := 0
	number := 0
	for index, pool := range *pools {
		go func(id uint32, addr *Pool) {
			zones, err := mds.GMdsClient.ListPoolZone(id)
			ret <- baserpc.RpcResult{
				Key:    addr,
				Result: zones,
				Err:    err,
			}
		}(pool.physicalPoolId, &(*pools)[index])
	}
	for res := range ret {
		if res.Err != nil {
			return res.Err
		}
		for _, z := range res.Result.([]common.Zone) {
			var zone Zone
			zone.Id = z.Id
			zone.Name = z.Name
			res.Key.(*Pool).Zones = append(res.Key.(*Pool).Zones, zone)
			number = number + 1
		}
		count = count + 1
		if count >= size {
			break
		}
	}
	return listZoneServer(pools, number)
}

func ListTopology() (interface{}, error) {
	var result []Pool
	logicalPools, err := mds.GMdsClient.ListLogicalPool()
	if err != nil {
		return nil, err
	}
	for _, lp := range logicalPools.([]common.LogicalPool) {
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
	return result, nil
}

func ListLogicalPool() (interface{}, error) {
	// get info from mds
	pools, err := mds.GMdsClient.ListLogicalPool()
	if err != nil {
		return nil, err
	}
	// TODO: get info from monitor
	return pools, nil
}
