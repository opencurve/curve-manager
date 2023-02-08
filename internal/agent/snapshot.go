package agent

import "github.com/opencurve/curve-manager/internal/snapshotclone"

func GetSnapshot(size, page uint32, uuid, user, fileName, status string) (interface{}, error) {
	return snapshotclone.GetSnapshot(size, page, uuid, user, fileName, status)
}
