/*
*  Copyright (c) 2023 NetEase Inc.
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
 */

/*
* Project: Curve-Manager
* Created Date: 2023-02-11
* Author: wanghai (SeanHai)
 */

package agent

import (
	"fmt"

	metricomm "github.com/opencurve/curve-manager/internal/metrics/common"
	"github.com/opencurve/curve-manager/internal/rpc/curvebs"
	bsrpc "github.com/opencurve/curve-manager/internal/rpc/curvebs"
)

const (
	ETCD_SERVICE                  = "etcd"
	MDS_SERVICE                   = "mds"
	SNAPSHOT_CLONE_SERVER_SERVICE = "snapshotcloneserver"

	RECYCLEBIN_DIR = "/RecycleBin"
)

type Server struct {
	Id           uint32              `json:"id" binding:"required"`
	Hostname     string              `json:"hostname" binding:"required"`
	InternalIp   string              `json:"internalIp" binding:"required"`
	InternalPort uint32              `json:"internalPort" binding:"required"`
	ExternalIp   string              `json:"externalIp" binding:"required"`
	ExternalPort uint32              `json:"externalPort" binding:"required"`
	ChunkServers []bsrpc.ChunkServer `json:"chunkservers" binding:"required"`
}

type Zone struct {
	Id      uint32   `json:"id" binding:"required"`
	Name    string   `json:"name" binding:"required"`
	Servers []Server `json:"servers" binding:"required"`
}

type Pool struct {
	Id             uint32 `json:"id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Type           string `json:"type" binding:"required"`
	Zones          []Zone `json:"zones" binding:"required"`
	physicalPoolId uint32
}

type Space struct {
	Total       uint64 `json:"total" binding:"required"`
	Alloc       uint64 `json:"alloc" binding:"required"`
	CanRecycled uint64 `json:"canRecycled" binding:"required"`
}

type PoolInfo struct {
	Id             uint32                  `json:"id" binding:"required"`
	Name           string                  `json:"name" binding:"required"`
	PhysicalPoolId uint32                  `json:"physicalPoolId" binding:"required"`
	Type           string                  `json:"type" binding:"required"`
	CreateTime     string                  `json:"createTime" binding:"required"`
	AllocateStatus string                  `json:"allocateStatus" binding:"required"`
	ScanEnable     bool                    `json:"scanEnable"`
	ServerNum      uint32                  `json:"serverNum" binding:"required"`
	ChunkServerNum uint32                  `json:"chunkServerNum" binding:"required"`
	CopysetNum     uint32                  `json:"copysetNum" binding:"required"`
	Space          Space                   `json:"space" binding:"required"`
	Performance    []metricomm.Performance `json:"performance" binding:"required"`
}

type VolumePoolInfo struct {
	Id    uint32 `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Alloc uint32 `json:"alloc" binding:"required"`
}
type VolumeInfo struct {
	Info        curvebs.FileInfo            `json:"info" binding:"required"`
	Pools       []VolumePoolInfo            `json:"pools"`
	Performance []metricomm.UserPerformance `json:"performance" binding:"required"`
}

type VersionNum struct {
	Version string `json:"version"`
	Number  int    `json:"number"`
}

type ChunkServerStatus struct {
	TotalNum  int          `json:"totalNum"`
	OnlineNum int          `json:"onlineNum"`
	Versions  []VersionNum `json:"versions"`
}

type CopysetNum struct {
	Total     uint32 `json:"total" binding:"required"`
	Unhealthy uint32 `json:"unhealthy" binding:"required"`
}

type ClusterStatus struct {
	Healthy    bool       `json:"healthy" binding:"required"`
	PoolNum    uint32     `json:"poolNum" binding:"required"`
	CopysetNum CopysetNum `json:"copysetNum" binding:"required"`
}

type HostInfo struct {
	HostName    string            `json:"hostName" binding:"required"`
	IP          string            `json:"ip" binding:"required"`
	Machine     string            `json:"machine" binding:"required"`
	Release     string            `json:"kernel-release" binding:"required"`
	Version     string            `json:"kernel-version" binding:"required"`
	System      string            `json:"operating-system" binding:"required"`
	CPUCores    metricomm.CPUInfo `json:"cpuCores" binding:"required"`
	DiskNum     uint32            `json:"diskNum" binding:"required"`
	MemoryTotal uint64            `json:"memory" binding:"required"`
}

type NetWorkTraffic struct {
	NetWorkReceive  map[string][]metricomm.RangeMetricItem `json:"receive" binding:"required"`
	NetWorkTransmit map[string][]metricomm.RangeMetricItem `json:"transmit" binding:"required"`
}

type HostPerformance struct {
	CPUUtilization  []metricomm.RangeMetricItem        `json:"cpuUtilization" binding:"required"`
	MemUtilization  []metricomm.RangeMetricItem        `json:"memUtilization" binding:"required"`
	DiskPerformance map[string][]metricomm.Performance `json:"diskPerformance" binding:"required"`
	NetWorkTraffic  NetWorkTraffic                     `json:"networkTraffic" binding:"required"`
}

type DiskInfo struct {
	HostName   string `json:"hostName" binding:"required"`
	Device     string `json:"device" binding:"required"`
	MountPoint string `json:"mountPoint"`
	FileSystem string `json:"fileSystem"`
	DiskType   string `json:"diskType"`
	SpaceTotal uint32 `json:"spaceTotal"`
	SpaceAvail uint32 `json:"spaceAvail"`
}

func getInstanceByHostName(hostname string) (string, error) {
	if hostname == "" {
		return "", nil
	}
	baseInfo, err := metricomm.GetHostsInfo()
	if err != nil {
		return "", err
	}
	for k, info := range baseInfo {
		if info.HostName == hostname {
			return k, nil
		}
	}
	return "", fmt.Errorf("hostname not exist, hostname = %s", hostname)
}

func getHostNameByInstance(instances []string) (map[string]string, error) {
	baseInfo, err := metricomm.GetHostsInfo()
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	for inst, info := range baseInfo {
		ret[inst] = info.HostName
	}
	return ret, nil
}
