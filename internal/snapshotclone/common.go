package snapshotclone

const (
	SERVICE_ROUTER         = "/SnapshotCloneService"
	SNAPSHOT_CLONE_VERSION = "0.0.6"

	ERROR_CODE_SUCCESS = "0"
)

// struct received from snapshot clone server
type SnapShotInfo struct {
	Code       string     `json:"Code" binding:"required"`
	Message    string     `json:"Message" binding:"required"`
	RequestId  string     `json:"RequestId" binding:"required"`
	Snapshots  []snapshot `json:"Snapshots"`
	TotalCount uint64     `json:"TotalCount"`
}

type snapshot struct {
	UUID       string `json:"UUID" binding:"required"`
	User       string `json:"User" binding:"required"`
	File       string `json:"File" binding:"required"`
	Name       string `json:"Name" binding:"required"`
	SeqNum     uint64 `json:"SeqNum" binding:"required"`
	Ctime      uint64 `json:"Time" binding:"required"`
	FileLength uint64 `json:"FileLength" binding:"required"`
	Status     int    `json:"Status" binding:"required"`
	Progress   uint32 `json:"Progress" binding:"required"`
}

// struct return to uplayer
type Snapshot struct {
	UUID       string `json:"uuid" binding:"required"`
	User       string `json:"user" binding:"required"`
	File       string `json:"file" binding:"required"`
	Name       string `json:"name" binding:"required"`
	SeqNum     uint64 `json:"seqNum" binding:"required"`
	Ctime      string `json:"time" binding:"required"`
	FileLength uint64 `json:"length" binding:"required"`
	Status     string `json:"status" binding:"required"`
	Progress   string `json:"Progress" binding:"required"`
}
