package common

type PhysicalPool struct {
	Id   uint32 `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc"`
}

type LogicalPool struct {
	Id               uint32 `json:"id" binding:"required"`
	Name             string `json:"name" binding:"required"`
	PhysicalPoolId   uint32 `json:"physicalPoolId" binding:"required"`
	Type             string `json:"type" binding:"required"`
	CreateTime       string `json:"createTime" binding:"required"`
	AllocateStatus   string `json:"allocateStatus" binding:"required"`
	ScanEnable       bool   `json:"scanEnable"`
	TotalSpace       string `json:"totalSpace" binding:"required"`
	UsedSpace        string `json:"UsedSpace" binding:"required"`
}
