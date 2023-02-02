package snapshotclone

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/opencurve/curve-manager/internal/common"
)

const (
	ACTION_GET_SNAPSHOT_LIST = "GetFileSnapshotList"

	// snapshot status
	STATUS_DONE           = "done"
	STATUS_IN_PROCESS     = "in-process"
	STATUS_DELETING       = "deleting"
	STATUS_ERROR_DELETING = "errorDeleting"
	STATUS_CANCELING      = "canceling"
	STATUS_ERROR          = "error"
	STATUS_NOT_SUPPORT    = "not-support"
)

func getNumericStatus(status string) int {
	switch status {
	case STATUS_DONE:
		return 0
	case STATUS_IN_PROCESS:
		return 1
	case STATUS_DELETING:
		return 2
	case STATUS_ERROR_DELETING:
		return 3
	case STATUS_CANCELING:
		return 4
	case STATUS_ERROR:
		return 5
	default:
		return -1
	}
}

func getStrStatus(status int) string {
	switch status {
	case 0:
		return STATUS_DONE
	case 1:
		return STATUS_IN_PROCESS
	case 2:
		return STATUS_DELETING
	case 3:
		return STATUS_ERROR_DELETING
	case 4:
		return STATUS_CANCELING
	case 5:
		return STATUS_ERROR
	default:
		return STATUS_NOT_SUPPORT
	}
}

func transferSnapshotInfo(in *[]snapshot, out *[]Snapshot) {
	for _, info := range *in {
		var item Snapshot
		item.UUID = info.UUID
		item.Name = info.Name
		item.User = info.User
		item.File = info.File
		item.SeqNum = info.SeqNum
		item.Ctime = time.Unix(int64(info.Ctime/1000000), 0).Format(common.TIME_FORMAT)
		item.FileLength = info.FileLength / common.GB
		item.Status = getStrStatus(info.Status)
		item.Progress = fmt.Sprintf("%d%%", info.Progress)
		*out = append(*out, item)
	}
}

func GetSnapshot(size, page uint32, uuid, user, fileName, status string) ([]Snapshot, error) {
	var snapshotInfo SnapShotInfo
	params := fmt.Sprintf("Action=%s&Version=%s&Limit=%d&Offset=%d",
	ACTION_GET_SNAPSHOT_LIST, SNAPSHOT_CLONE_VERSION, size, (page-1)*size)
	if uuid != "" {
		params = fmt.Sprintf("%s&UUID=%s", params, uuid)
	}
	if user != "" {
		params = fmt.Sprintf("%s&User=%s", params, user)
	}
	if fileName != "" {
		params = fmt.Sprintf("%s&File=%s", params, fileName)
	}
	if status != "" {
		s := getNumericStatus(status)
		if s < 0 {
			return nil, fmt.Errorf("status not support")
		}
		params = fmt.Sprintf("%s&Status=%d", params, s)
	}

	resp, err := GSnapshotCloneClient.sendHttp2SnapshotClone(params)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(resp), &snapshotInfo)
	if err != nil {
		return nil, err
	}

	if snapshotInfo.Code != ERROR_CODE_SUCCESS {
		return nil, fmt.Errorf(snapshotInfo.Message)
	}

	var snapshots []Snapshot
	transferSnapshotInfo(&snapshotInfo.Snapshots, &snapshots)
	return snapshots, nil
}
