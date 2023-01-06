package agent

import "github.com/opencurve/curve-manager/internal/rpc/curvebs/common"

type Server struct {
	Id           uint32               `json:"id" binding:"required"`
	Hostname     string               `json:"hostname" binding:"required"`
	InternalIp   string               `json:"internalIp" binding:"required"`
	InternalPort uint32               `json:"internalPort" binding:"required"`
	ExternalIp   string               `json:"externalIp" binding:"required"`
	ExternalPort uint32               `json:"externalPort" binding:"required"`
	ChunkServers []common.ChunkServer `json:"chunkservers" binding:"required"`
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
