/*
*  Copyright (c) 2023 NetEase Inc.
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
 */

/*
* Project: Curve-Manager
* Created Date: 2023-03-22
* Author: wanghai (SeanHai)
 */

package agent

import (
	"fmt"
	"time"

	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/storage"
	"github.com/opencurve/pigeon"
)

const (
	ALERT_REQUEST_ID = "alert"

	MODULE_CLUSTER = "cluster"
	MODULE_SPACE   = "space"

	CAPACITY_ALERT_INTERVAL_SEC  = 1 * 60 * 60
	CAPACITY_ALERT_TRIGGER_TIMES = 1
	CAPACITY_ALERT_LIMIT_PERCENT = 80

	SERVICE_ALERT_INTERVAL_SEC  = 1 * 60
	SERVICE_ALERT_TRIGGER_TIMES = 1

	CLUSTER_ALERT_INTERVAL_SEC  = 1 * 60
	CLUSTER_ALERT_TRIGGER_TIMES = 3
)

var (
	capacityAlertClient *capacityAlert
	serviceAlertClient  *serviceAlert
	clusterAlertClient  *clusterAlert
)

type checkOption struct {
	checkIntervalSec uint32
	triggeredTimes   uint32
}

type Alert interface {
	check(r *pigeon.Request)
}

type capacityAlert struct {
	opt          checkOption
	limitPercent uint32
	times        uint32
}

type serviceAlert struct {
	opt      checkOption
	services map[string]uint32
}

type clusterAlert struct {
	opt   checkOption
	times uint32
}

func initAlerts(logger *pigeon.Logger) {
	capacityAlertClient = &capacityAlert{
		times:        0,
		limitPercent: CAPACITY_ALERT_LIMIT_PERCENT,
		opt: checkOption{
			checkIntervalSec: CAPACITY_ALERT_INTERVAL_SEC,
			triggeredTimes:   CAPACITY_ALERT_TRIGGER_TIMES,
		},
	}
	go capacityAlertClient.check(logger)

	serviceAlertClient = &serviceAlert{
		opt: checkOption{
			checkIntervalSec: SERVICE_ALERT_INTERVAL_SEC,
			triggeredTimes:   SERVICE_ALERT_TRIGGER_TIMES,
		},
		services: map[string]uint32{
			SERVICE_ETCD:                  0,
			SERVICE_MDS:                   0,
			SERVICE_CHUNKSERVER:           0,
			SERVICE_SNAPSHOT_CLONE_SERVER: 0,
		},
	}
	go serviceAlertClient.check(logger)

	clusterAlertClient = &clusterAlert{
		times: 0,
		opt: checkOption{
			checkIntervalSec: CLUSTER_ALERT_INTERVAL_SEC,
			triggeredTimes:   CLUSTER_ALERT_TRIGGER_TIMES,
		},
	}
	go clusterAlertClient.check(logger)
}

func recordAlert(level int, module string, duration uint32, summary string) error {
	now := time.Now().UnixMilli()
	return storage.AddSystemAlert(&storage.SystemAlert{
		TimeMs:      now,
		Level:       level,
		Module:      module,
		DurationSec: duration,
		Summary:     summary,
	})
}

func GetUnreadSysAlertNum(r *pigeon.Request) (int64, errno.Errno) {
	token := r.HeadersIn[comm.HEADER_AUTH_TOKEN]
	user := storage.GetLoginUserByToken(token)
	if user == "" {
		r.Logger().Error("GetUnreadSysAlertNum get user by token failed",
			pigeon.Field("token", token),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return 0, errno.GET_USER_FAILED
	}
	maxId, err := storage.GetLastAlertId()
	if err != nil {
		r.Logger().Error("GetLastAlertId failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return 0, errno.GET_LAST_ALERT_ID_FAILED
	}
	readId, err := storage.GetReadAlertId(user)
	if err != nil {
		r.Logger().Error("GetReadAlertId failed",
			pigeon.Field("user", user),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return 0, errno.GET_READ_ALERT_ID_FAILED
	}
	if readId == -1 {
		readId = 0
		err = storage.AddReadAlertId(user)
		if err != nil {
			r.Logger().Error("AddReadAlertId failed",
				pigeon.Field("user", user),
				pigeon.Field("error", err),
				pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
			return 0, errno.ADD_READ_ALERT_ID_FAILED
		}
	}
	return maxId-readId, errno.OK
}

func UpdateReadSysAlertId(r *pigeon.Request, id int64) errno.Errno {
	token := r.HeadersIn[comm.HEADER_AUTH_TOKEN]
	user := storage.GetLoginUserByToken(token)
	err := storage.UpdateReadAlertId(id, user)
	if err != nil {
		r.Logger().Error("UpdateReadSysAlertId failed",
			pigeon.Field("user", user),
			pigeon.Field("id", id),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return errno.UPDATE_UNREAD_ALERT_ID_FAILED
	}
	return errno.OK
}

func GetSysAlert(r *pigeon.Request, start, end int64, page, size uint32, filter string) (interface{}, errno.Errno) {
	if start == 0 && end == 0 {
		end = time.Now().UnixMilli()
	}
	info, err := storage.GetSystemAlert(start, end, size, (page-1)*size, filter)
	if err != nil {
		r.Logger().Error("GetAlert failed",
			pigeon.Field("start", start),
			pigeon.Field("end", end),
			pigeon.Field("filter", filter),
			pigeon.Field("page", page),
			pigeon.Field("size", size),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_SYSTEM_ALERT_FAILED
	}
	return info, errno.OK
}

func (alert *capacityAlert) check(logger *pigeon.Logger) {
	timer := time.NewTimer(time.Duration(alert.opt.checkIntervalSec) * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			space, err := GetClusterSpace(logger, ALERT_REQUEST_ID)
			if err == errno.OK {
				percent := space.(Space).Alloc * 100 / space.(Space).Total
				if percent >= uint64(alert.limitPercent) {
					alert.times++
					if alert.times >= alert.opt.triggeredTimes {
						e := recordAlert(storage.ALERT_WARNING, MODULE_SPACE, alert.opt.triggeredTimes*alert.opt.checkIntervalSec,
							fmt.Sprintf("Cluster space have alloced %d%%(alert trigger is %d%%)", percent, alert.limitPercent))
						if e != nil {
							logger.Error("record space alert falied",
								pigeon.Field("error", e))
						}
						alert.times = 0
					}
				}
			}
			timer.Reset(time.Duration(alert.opt.checkIntervalSec) * time.Second)
		}
	}
}

func (alert *serviceAlert) check(logger *pigeon.Logger) {
	timer := time.NewTimer(time.Duration(alert.opt.checkIntervalSec) * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			for name, times := range alert.services {
				if name == SERVICE_CHUNKSERVER {
					csStatus, err := GetChunkServerStatus(logger, ALERT_REQUEST_ID)
					if err == errno.OK {
						offline := csStatus.(*ChunkServerStatus).TotalNum - csStatus.(*ChunkServerStatus).OnlineNum
						if offline > 0 {
							times++
							if times >= alert.opt.triggeredTimes {
								e := recordAlert(storage.ALERT_CRITICAL, name, alert.opt.checkIntervalSec*alert.opt.triggeredTimes,
									fmt.Sprintf("chunkserver offline number = %d", offline))
								if e != nil {
									logger.Error("record service alert info failed",
										pigeon.Field("error", e))
								}
								times = 0
							}
							alert.services[name] = times
						}
					}
				} else {
					result := checkServiceHealthy(name)
					if !result.Result.(serviceStatus).healthy {
						times++
						if times >= alert.opt.triggeredTimes {
							e := recordAlert(storage.ALERT_CRITICAL, name, alert.opt.checkIntervalSec*alert.opt.triggeredTimes,
								result.Result.(serviceStatus).detail)
							if e != nil {
								logger.Error("record service alert info failed",
									pigeon.Field("error", e))
							}
							times = 0
						}
						alert.services[name] = times
					}
				}
			}
			timer.Reset(time.Duration(alert.opt.checkIntervalSec) * time.Second)
		}
	}
}

func (alert *clusterAlert) check(logger *pigeon.Logger) {
	timer := time.NewTimer(time.Duration(alert.opt.checkIntervalSec) * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			status := GetClusterStatus(logger, ALERT_REQUEST_ID)
			if !status.(ClusterStatus).Healthy {
				alert.times++
				if alert.times >= alert.opt.triggeredTimes {
					summary := "cluster is not healthy"
					unhealthy := status.(ClusterStatus).CopysetNum.Unhealthy
					if unhealthy > 0 {
						summary = fmt.Sprintf("%s, unhealthy copysets number is %d", summary, unhealthy)
					}
					e := recordAlert(storage.ALERT_WARNING, MODULE_CLUSTER, alert.opt.checkIntervalSec*alert.opt.triggeredTimes,
						summary)
					if e != nil {
						logger.Error("record cluster alert info failed",
							pigeon.Field("error", e))
					}
					alert.times = 0
				}
			}
			timer.Reset(time.Duration(alert.opt.checkIntervalSec) * time.Second)
		}
	}
}
