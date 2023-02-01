package agent

import (
	bsmetricomm "github.com/opencurve/curve-manager/internal/metrics/common"
	bsrpc "github.com/opencurve/curve-manager/internal/rpc/curvebs"
)

const (
	ETCD_SERVICE                  = "etcd"
	MDS_SERVICE                   = "mds"
	SNAPSHOT_CLONE_SERVER_SERVICE = "snapshotcloneserver"
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
	Id             uint32               `json:"id" binding:"required"`
	Name           string               `json:"name" binding:"required"`
	PhysicalPoolId uint32               `json:"physicalPoolId" binding:"required"`
	Type           string               `json:"type" binding:"required"`
	CreateTime     string               `json:"createTime" binding:"required"`
	AllocateStatus string               `json:"allocateStatus" binding:"required"`
	ScanEnable     bool                 `json:"scanEnable"`
	ServerNum      uint32               `json:"serverNum" binding:"required"`
	ChunkServerNum uint32               `json:"chunkServerNum" binding:"required"`
	CopysetNum     uint32               `json:"copysetNum" binding:"required"`
	Space          Space                `json:"space" binding:"required"`
	Performance    []bsmetricomm.Performance `json:"performance" binding:"required"`
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
