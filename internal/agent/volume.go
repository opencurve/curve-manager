package agent

import (
	"fmt"
	"sort"
	"time"

	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/metrics/bsmetric"
	bsrpc "github.com/opencurve/curve-manager/internal/rpc/curvebs"
)

const (
	ORDER_BY_ID              = "id"
	ORDER_BY_CTIME           = "ctime"
	ORDER_BY_LENGTH          = "length"
	ORDER_DIRECTION_INCREASE = 1
	ORDER_DIRECTION_DECREASE = -1
)

func sortFile(files []bsrpc.FileInfo, orderKey string, direction int) {
	sort.Slice(files, func(i, j int) bool {
		switch orderKey {
		case ORDER_BY_CTIME:
			itime, _ := time.Parse(common.TIME_FORMAT, files[i].Ctime)
			jtime, _ := time.Parse(common.TIME_FORMAT, files[j].Ctime)
			if direction == ORDER_DIRECTION_INCREASE {
				return itime.Unix() < jtime.Unix()
			}
			return itime.Unix() > jtime.Unix()
		case ORDER_BY_LENGTH:
			if direction == ORDER_DIRECTION_INCREASE {
				return files[i].Length < files[j].Length
			}
			return files[i].Length > files[j].Length
		}
		if direction == ORDER_DIRECTION_INCREASE {
			return files[i].Id < files[j].Id
		}
		return files[i].Id > files[j].Id
	})
}

func ListVolume(size, page uint32, path, key string, direction int) (interface{}, error) {
	// get root auth info
	authInfo, err := bsmetric.GetAuthInfoOfRoot()
	if err != "" {
		return nil, fmt.Errorf(err)
	}

	// create signature
	date := time.Now().UnixMicro()
	str2sig := common.GetString2Signature(date, authInfo.UserName)
	sig := common.CalcString2Signature(str2sig, authInfo.PassWord)

	fileInfos, e := bsrpc.GMdsClient.ListDir(path, authInfo.UserName, sig, uint64(date))
	if e != nil {
		return nil, fmt.Errorf("ListDir failed, %s", e)
	}

	if len(fileInfos) == 0 {
		return nil, nil
	}

	sortFile(fileInfos, key, direction)
	length := uint32(len(fileInfos))
	start := (page - 1) * size
	var end uint32
	if page*size > length {
		end = length
	} else {
		end = page * size
	}
	return fileInfos[start:end], nil
}
