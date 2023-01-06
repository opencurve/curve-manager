package common

type PhysicalPool struct {
	Id   uint32 `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc"`
}

type LogicalPool struct {
	Id             uint32 `json:"id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	PhysicalPoolId uint32 `json:"physicalPoolId" binding:"required"`
	Type           string `json:"type" binding:"required"`
	CreateTime     string `json:"createTime" binding:"required"`
	AllocateStatus string `json:"allocateStatus" binding:"required"`
	ScanEnable     bool   `json:"scanEnable"`
}

type Zone struct {
	Id               uint32 `json:"id" binding:"required"`
	Name             string `json:"name" binding:"required"`
	PhysicalPoolId   uint32 `json:"physicalPoolId" binding:"required"`
	PhysicalPoolName string `json:"physicalName" binding:"required"`
	Desc             string `json:"desc"`
}

type Server struct {
	Id               uint32 `json:"id" binding:"required"`
	HostName         string `json:"hostName" binding:"required"`
	InternalIp       string `json:"internalIp" binding:"required"`
	InternalPort     uint32 `json:"internalPort" binding:"required"`
	ExternalIp       string `json:"externalIp" binding:"required"`
	ExternalPort     uint32 `json:"externalPort" binding:"required"`
	ZoneId           uint32 `json:"zoneId" binding:"required"`
	ZoneName         string `json:"zoneName" binding:"required"`
	PhysicalPoolId   uint32 `json:"physicalPoolId" binding:"required"`
	PhysicalPoolName string `json:"physicalName" binding:"required"`
	Desc             string `json:"desc"`
}

type ChunkServer struct {
	Id           uint32 `json:"id" binding:"required"`
	DiskType     string `json:"diskType" binding:"required"`
	HostIp       string `json:"hostIp" binding:"required"`
	Port         uint32 `json:"port" binding:"required"`
	Status       string `json:"status" binding:"required"`
	DiskStatus   string `json:"diskStatus" binding:"required"`
	OnlineStatus string `json:"onlineStatus" binding:"required"`
	MountPoint   string `json:"mountPoint" binding:"required"`
	DiskCapacity string `json:"diskCapacity" binding:"required"`
	DiskUsed     string `json:"diskUsed" binding:"required"`
	ExternalIp   string `json:"externalIp"`
}
