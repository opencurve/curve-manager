package common

type PhysicalPool struct {
	Id uint32 `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc" binding:"required"`
}