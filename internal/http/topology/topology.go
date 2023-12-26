package topology

import (
	"github.com/opencurve/curve-manager/internal/http/common"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type LogicalPoolType int32

const (
	PAGEFILE     LogicalPoolType = 0
	APPENDFILE   LogicalPoolType = 1
	APPENDECFILE LogicalPoolType = 2
)

type AllocateStatus int32

const (
	ALLOW AllocateStatus = 0
	DENY  AllocateStatus = 1
)

type ChunkServerStatus int32

const (
	READWRITE ChunkServerStatus = 0
	PENDDING  ChunkServerStatus = 1
	RETIRED   ChunkServerStatus = 2
)

type DiskState int32

const (
	DISKNORMAL DiskState = 0
	DISKERROR  DiskState = 1
)

type OnlineState int32

const (
	ONLINE   OnlineState = 0
	OFFLINE  OnlineState = 1
	UNSTABLE OnlineState = 2
)

type ListPhysicalPoolResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode        *int32              `protobuf:"zigzag32,1,req,name=statusCode" json:"statusCode,omitempty"`
	PhysicalPoolInfos []*PhysicalPoolInfo `protobuf:"bytes,2,rep,name=physicalPoolInfos" json:"physicalPoolInfos,omitempty"`
}

type PhysicalPoolInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PhysicalPoolID   *uint32 `protobuf:"varint,1,req,name=physicalPoolID" json:"physicalPoolID,omitempty"`
	PhysicalPoolName *string `protobuf:"bytes,2,req,name=physicalPoolName" json:"physicalPoolName,omitempty"`
	Desc             *string `protobuf:"bytes,3,opt,name=desc" json:"desc,omitempty"`
}

func (x *ListPhysicalPoolResponse) GetStatusCode() int32 {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return 0
}

func (x *ListPhysicalPoolResponse) GetPhysicalPoolInfos() []*PhysicalPoolInfo {
	if x != nil {
		return x.PhysicalPoolInfos
	}
	return nil
}

func (x *PhysicalPoolInfo) GetPhysicalPoolID() uint32 {
	if x != nil && x.PhysicalPoolID != nil {
		return *x.PhysicalPoolID
	}
	return 0
}

func (x *PhysicalPoolInfo) GetPhysicalPoolName() string {
	if x != nil && x.PhysicalPoolName != nil {
		return *x.PhysicalPoolName
	}
	return ""
}

func (x *PhysicalPoolInfo) GetDesc() string {
	if x != nil && x.Desc != nil {
		return *x.Desc
	}
	return ""
}

type ListLogicalPoolResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode       *int32             `protobuf:"zigzag32,1,req,name=statusCode" json:"statusCode,omitempty"`
	LogicalPoolInfos []*LogicalPoolInfo `protobuf:"bytes,2,rep,name=logicalPoolInfos" json:"logicalPoolInfos,omitempty"`
}

func (x *ListLogicalPoolResponse) GetStatusCode() int32 {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return 0
}

func (x *ListLogicalPoolResponse) GetLogicalPoolInfos() []*LogicalPoolInfo {
	if x != nil {
		return x.LogicalPoolInfos
	}
	return nil
}

type LogicalPoolInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LogicalPoolID                *uint32          `protobuf:"varint,1,req,name=logicalPoolID" json:"logicalPoolID,omitempty"`
	LogicalPoolName              *string          `protobuf:"bytes,2,req,name=logicalPoolName" json:"logicalPoolName,omitempty"`
	PhysicalPoolID               *uint32          `protobuf:"varint,3,req,name=physicalPoolID" json:"physicalPoolID,omitempty"`
	Type                         *LogicalPoolType `protobuf:"varint,4,req,name=type,enum=curve.mds.topology.LogicalPoolType" json:"type,omitempty"`
	CreateTime                   *uint64          `protobuf:"varint,5,req,name=createTime" json:"createTime,omitempty"`
	RedundanceAndPlaceMentPolicy []byte           `protobuf:"bytes,6,req,name=redundanceAndPlaceMentPolicy" json:"redundanceAndPlaceMentPolicy,omitempty"` //json body
	UserPolicy                   []byte           `protobuf:"bytes,7,req,name=userPolicy" json:"userPolicy,omitempty"`                                     //json body
	AllocateStatus               *AllocateStatus  `protobuf:"varint,8,req,name=allocateStatus,enum=curve.mds.topology.AllocateStatus" json:"allocateStatus,omitempty"`
	ScanEnable                   *bool            `protobuf:"varint,9,opt,name=scanEnable" json:"scanEnable,omitempty"`
}

func (x *LogicalPoolInfo) GetLogicalPoolID() uint32 {
	if x != nil && x.LogicalPoolID != nil {
		return *x.LogicalPoolID
	}
	return 0
}

func (x *LogicalPoolInfo) GetLogicalPoolName() string {
	if x != nil && x.LogicalPoolName != nil {
		return *x.LogicalPoolName
	}
	return ""
}

func (x *LogicalPoolInfo) GetPhysicalPoolID() uint32 {
	if x != nil && x.PhysicalPoolID != nil {
		return *x.PhysicalPoolID
	}
	return 0
}

func (x *LogicalPoolInfo) GetType() LogicalPoolType {
	if x != nil && x.Type != nil {
		return *x.Type
	}
	return PAGEFILE
}

func (x *LogicalPoolInfo) GetCreateTime() uint64 {
	if x != nil && x.CreateTime != nil {
		return *x.CreateTime
	}
	return 0
}

func (x *LogicalPoolInfo) GetRedundanceAndPlaceMentPolicy() []byte {
	if x != nil {
		return x.RedundanceAndPlaceMentPolicy
	}
	return nil
}

func (x *LogicalPoolInfo) GetUserPolicy() []byte {
	if x != nil {
		return x.UserPolicy
	}
	return nil
}

func (x *LogicalPoolInfo) GetAllocateStatus() AllocateStatus {
	if x != nil && x.AllocateStatus != nil {
		return *x.AllocateStatus
	}
	return ALLOW
}

func (x *LogicalPoolInfo) GetScanEnable() bool {
	if x != nil && x.ScanEnable != nil {
		return *x.ScanEnable
	}
	return false
}

type GetLogicalPoolResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode      *int32           `protobuf:"zigzag32,1,req,name=statusCode" json:"statusCode,omitempty"`
	LogicalPoolInfo *LogicalPoolInfo `protobuf:"bytes,2,opt,name=logicalPoolInfo" json:"logicalPoolInfo,omitempty"`
}

func (x *GetLogicalPoolResponse) GetStatusCode() int32 {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return 0
}

func (x *GetLogicalPoolResponse) GetLogicalPoolInfo() *LogicalPoolInfo {
	if x != nil {
		return x.LogicalPoolInfo
	}
	return nil
}

type ListPoolZoneResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode *int32      `protobuf:"zigzag32,1,req,name=statusCode" json:"statusCode,omitempty"`
	Zones      []*ZoneInfo `protobuf:"bytes,2,rep,name=zones" json:"zones,omitempty"`
}

func (x *ListPoolZoneResponse) GetStatusCode() int32 {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return 0
}

func (x *ListPoolZoneResponse) GetZones() []*ZoneInfo {
	if x != nil {
		return x.Zones
	}
	return nil
}

type ZoneInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ZoneID           *uint32 `protobuf:"varint,1,req,name=zoneID" json:"zoneID,omitempty"`
	ZoneName         *string `protobuf:"bytes,2,req,name=zoneName" json:"zoneName,omitempty"`
	PhysicalPoolID   *uint32 `protobuf:"varint,3,req,name=physicalPoolID" json:"physicalPoolID,omitempty"`
	PhysicalPoolName *string `protobuf:"bytes,4,req,name=physicalPoolName" json:"physicalPoolName,omitempty"`
	Desc             *string `protobuf:"bytes,5,opt,name=desc" json:"desc,omitempty"`
}

func (x *ZoneInfo) GetZoneID() uint32 {
	if x != nil && x.ZoneID != nil {
		return *x.ZoneID
	}
	return 0
}

func (x *ZoneInfo) GetZoneName() string {
	if x != nil && x.ZoneName != nil {
		return *x.ZoneName
	}
	return ""
}

func (x *ZoneInfo) GetPhysicalPoolID() uint32 {
	if x != nil && x.PhysicalPoolID != nil {
		return *x.PhysicalPoolID
	}
	return 0
}

func (x *ZoneInfo) GetPhysicalPoolName() string {
	if x != nil && x.PhysicalPoolName != nil {
		return *x.PhysicalPoolName
	}
	return ""
}

func (x *ZoneInfo) GetDesc() string {
	if x != nil && x.Desc != nil {
		return *x.Desc
	}
	return ""
}

type ListZoneServerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode *int32        `protobuf:"zigzag32,1,req,name=statusCode" json:"statusCode,omitempty"`
	ServerInfo []*ServerInfo `protobuf:"bytes,2,rep,name=serverInfo" json:"serverInfo,omitempty"`
}

func (x *ListZoneServerResponse) GetStatusCode() int32 {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return 0
}

func (x *ListZoneServerResponse) GetServerInfo() []*ServerInfo {
	if x != nil {
		return x.ServerInfo
	}
	return nil
}

type ServerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServerID         *uint32 `protobuf:"varint,1,req,name=serverID" json:"serverID,omitempty"`
	HostName         *string `protobuf:"bytes,2,req,name=hostName" json:"hostName,omitempty"`
	InternalIp       *string `protobuf:"bytes,3,req,name=internalIp" json:"internalIp,omitempty"`
	InternalPort     *uint32 `protobuf:"varint,4,req,name=internalPort" json:"internalPort,omitempty"`
	ExternalIp       *string `protobuf:"bytes,5,req,name=externalIp" json:"externalIp,omitempty"`
	ExternalPort     *uint32 `protobuf:"varint,6,req,name=externalPort" json:"externalPort,omitempty"`
	ZoneID           *uint32 `protobuf:"varint,7,req,name=zoneID" json:"zoneID,omitempty"`
	ZoneName         *string `protobuf:"bytes,8,req,name=zoneName" json:"zoneName,omitempty"`
	PhysicalPoolID   *uint32 `protobuf:"varint,9,req,name=physicalPoolID" json:"physicalPoolID,omitempty"`
	PhysicalPoolName *string `protobuf:"bytes,10,req,name=physicalPoolName" json:"physicalPoolName,omitempty"`
	Desc             *string `protobuf:"bytes,11,req,name=desc" json:"desc,omitempty"`
}

func (x *ServerInfo) GetServerID() uint32 {
	if x != nil && x.ServerID != nil {
		return *x.ServerID
	}
	return 0
}

func (x *ServerInfo) GetHostName() string {
	if x != nil && x.HostName != nil {
		return *x.HostName
	}
	return ""
}

func (x *ServerInfo) GetInternalIp() string {
	if x != nil && x.InternalIp != nil {
		return *x.InternalIp
	}
	return ""
}

func (x *ServerInfo) GetInternalPort() uint32 {
	if x != nil && x.InternalPort != nil {
		return *x.InternalPort
	}
	return 0
}

func (x *ServerInfo) GetExternalIp() string {
	if x != nil && x.ExternalIp != nil {
		return *x.ExternalIp
	}
	return ""
}

func (x *ServerInfo) GetExternalPort() uint32 {
	if x != nil && x.ExternalPort != nil {
		return *x.ExternalPort
	}
	return 0
}

func (x *ServerInfo) GetZoneID() uint32 {
	if x != nil && x.ZoneID != nil {
		return *x.ZoneID
	}
	return 0
}

func (x *ServerInfo) GetZoneName() string {
	if x != nil && x.ZoneName != nil {
		return *x.ZoneName
	}
	return ""
}

func (x *ServerInfo) GetPhysicalPoolID() uint32 {
	if x != nil && x.PhysicalPoolID != nil {
		return *x.PhysicalPoolID
	}
	return 0
}

func (x *ServerInfo) GetPhysicalPoolName() string {
	if x != nil && x.PhysicalPoolName != nil {
		return *x.PhysicalPoolName
	}
	return ""
}

func (x *ServerInfo) GetDesc() string {
	if x != nil && x.Desc != nil {
		return *x.Desc
	}
	return ""
}

type ListChunkServerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode       *int32             `protobuf:"zigzag32,1,req,name=statusCode" json:"statusCode,omitempty"`
	ChunkServerInfos []*ChunkServerInfo `protobuf:"bytes,2,rep,name=chunkServerInfos" json:"chunkServerInfos,omitempty"`
}

func (x *ListChunkServerResponse) GetStatusCode() int32 {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return 0
}

func (x *ListChunkServerResponse) GetChunkServerInfos() []*ChunkServerInfo {
	if x != nil {
		return x.ChunkServerInfos
	}
	return nil
}

type ChunkServerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChunkServerID *uint32            `protobuf:"varint,1,req,name=chunkServerID" json:"chunkServerID,omitempty"`
	DiskType      *string            `protobuf:"bytes,2,req,name=diskType" json:"diskType,omitempty"`
	HostIp        *string            `protobuf:"bytes,3,req,name=hostIp" json:"hostIp,omitempty"`
	Port          *uint32            `protobuf:"varint,4,req,name=port" json:"port,omitempty"`
	Status        *ChunkServerStatus `protobuf:"varint,5,req,name=status,enum=curve.mds.topology.ChunkServerStatus" json:"status,omitempty"`
	DiskStatus    *DiskState         `protobuf:"varint,6,req,name=diskStatus,enum=curve.mds.topology.DiskState" json:"diskStatus,omitempty"`
	OnlineState   *OnlineState       `protobuf:"varint,7,req,name=onlineState,enum=curve.mds.topology.OnlineState" json:"onlineState,omitempty"`
	MountPoint    *string            `protobuf:"bytes,8,req,name=mountPoint" json:"mountPoint,omitempty"`
	DiskCapacity  *uint64            `protobuf:"varint,9,req,name=diskCapacity" json:"diskCapacity,omitempty"`
	DiskUsed      *uint64            `protobuf:"varint,10,req,name=diskUsed" json:"diskUsed,omitempty"`
	ExternalIp    *string            `protobuf:"bytes,11,opt,name=externalIp" json:"externalIp,omitempty"`
}

func (x *ChunkServerInfo) GetChunkServerID() uint32 {
	if x != nil && x.ChunkServerID != nil {
		return *x.ChunkServerID
	}
	return 0
}

func (x *ChunkServerInfo) GetDiskType() string {
	if x != nil && x.DiskType != nil {
		return *x.DiskType
	}
	return ""
}

func (x *ChunkServerInfo) GetHostIp() string {
	if x != nil && x.HostIp != nil {
		return *x.HostIp
	}
	return ""
}

func (x *ChunkServerInfo) GetPort() uint32 {
	if x != nil && x.Port != nil {
		return *x.Port
	}
	return 0
}

func (x *ChunkServerInfo) GetStatus() ChunkServerStatus {
	if x != nil && x.Status != nil {
		return *x.Status
	}
	return READWRITE
}

func (x *ChunkServerInfo) GetDiskStatus() DiskState {
	if x != nil && x.DiskStatus != nil {
		return *x.DiskStatus
	}
	return DISKNORMAL
}

func (x *ChunkServerInfo) GetOnlineState() OnlineState {
	if x != nil && x.OnlineState != nil {
		return *x.OnlineState
	}
	return ONLINE
}

func (x *ChunkServerInfo) GetMountPoint() string {
	if x != nil && x.MountPoint != nil {
		return *x.MountPoint
	}
	return ""
}

func (x *ChunkServerInfo) GetDiskCapacity() uint64 {
	if x != nil && x.DiskCapacity != nil {
		return *x.DiskCapacity
	}
	return 0
}

func (x *ChunkServerInfo) GetDiskUsed() uint64 {
	if x != nil && x.DiskUsed != nil {
		return *x.DiskUsed
	}
	return 0
}

func (x *ChunkServerInfo) GetExternalIp() string {
	if x != nil && x.ExternalIp != nil {
		return *x.ExternalIp
	}
	return ""
}

type GetChunkServerInClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode       *int32             `protobuf:"zigzag32,1,req,name=statusCode" json:"statusCode,omitempty"`
	ChunkServerInfos []*ChunkServerInfo `protobuf:"bytes,2,rep,name=chunkServerInfos" json:"chunkServerInfos,omitempty"`
}

func (x *GetChunkServerInClusterResponse) GetStatusCode() int32 {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return 0
}

func (x *GetChunkServerInClusterResponse) GetChunkServerInfos() []*ChunkServerInfo {
	if x != nil {
		return x.ChunkServerInfos
	}
	return nil
}

type GetCopySetsInChunkServerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode   *int32                `protobuf:"zigzag32,1,req,name=statusCode" json:"statusCode,omitempty"`
	CopysetInfos []*common.CopysetInfo `protobuf:"bytes,2,rep,name=copysetInfos" json:"copysetInfos,omitempty"`
}

func (x *GetCopySetsInChunkServerResponse) GetStatusCode() int32 {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return 0
}

func (x *GetCopySetsInChunkServerResponse) GetCopysetInfos() []*common.CopysetInfo {
	if x != nil {
		return x.CopysetInfos
	}
	return nil
}

type GetChunkServerListInCopySetsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode *int32               `protobuf:"zigzag32,1,req,name=statusCode" json:"statusCode,omitempty"`
	CsInfo     []*CopySetServerInfo `protobuf:"bytes,2,rep,name=csInfo" json:"csInfo,omitempty"`
}

func (x *GetChunkServerListInCopySetsResponse) GetStatusCode() int32 {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return 0
}

func (x *GetChunkServerListInCopySetsResponse) GetCsInfo() []*CopySetServerInfo {
	if x != nil {
		return x.CsInfo
	}
	return nil
}

type CopySetServerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CopysetId *uint32                       `protobuf:"varint,1,req,name=copysetId" json:"copysetId,omitempty"`
	CsLocs    []*common.ChunkServerLocation `protobuf:"bytes,2,rep,name=csLocs" json:"csLocs,omitempty"`
}

func (x *CopySetServerInfo) GetCopysetId() uint32 {
	if x != nil && x.CopysetId != nil {
		return *x.CopysetId
	}
	return 0
}

func (x *CopySetServerInfo) GetCsLocs() []*common.ChunkServerLocation {
	if x != nil {
		return x.CsLocs
	}
	return nil
}

type GetCopySetsInClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode   *int32                `protobuf:"zigzag32,1,req,name=statusCode" json:"statusCode,omitempty"`
	CopysetInfos []*common.CopysetInfo `protobuf:"bytes,2,rep,name=copysetInfos" json:"copysetInfos,omitempty"`
}

func (x *GetCopySetsInClusterResponse) GetStatusCode() int32 {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return 0
}

func (x *GetCopySetsInClusterResponse) GetCopysetInfos() []*common.CopysetInfo {
	if x != nil {
		return x.CopysetInfos
	}
	return nil
}
