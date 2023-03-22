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

	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/storage"
	"github.com/opencurve/pigeon"
)

const (
	MODULE_CLUSTER = "cluster"

	CAPACITY_ALERT_INTERVAL_SEC  = 12 * 60 * 60
	CAPACITY_ALERT_TRIGGER_TIMES = 1
	CAPACITY_ALERT_LIMIT_PERCENT = 80
)

var (
	capacityAlertClient *capacityAlert
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

func initAlerts(r *pigeon.Request) {
	capacityAlertClient = &capacityAlert{
		times:        0,
		limitPercent: CAPACITY_ALERT_LIMIT_PERCENT,
		opt: checkOption{
			checkIntervalSec: CAPACITY_ALERT_INTERVAL_SEC,
			triggeredTimes:   CAPACITY_ALERT_TRIGGER_TIMES,
		},
	}
	go capacityAlertClient.check(r)
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

func (alert *capacityAlert) check(r *pigeon.Request) {
	timer := time.NewTimer(time.Duration(alert.opt.checkIntervalSec) * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			space, err := GetClusterSpace(r)
			if err == errno.OK {
				percent := space.(Space).Alloc * 100 / space.(Space).Total
				if percent >= uint64(alert.limitPercent) {
					alert.times++
					if alert.times >= alert.opt.triggeredTimes {
						e := recordAlert(storage.ALERT_WARNING, MODULE_CLUSTER, alert.opt.triggeredTimes*alert.opt.checkIntervalSec,
							fmt.Sprintf("Cluster space have alloced %d(alert trigger is %d)", percent, alert.limitPercent))
						if e != nil {
							r.Logger().Error("record cluster space alert falied",
								pigeon.Field("error", e))
						}
					}
				}
			}
			timer.Reset(time.Duration(alert.opt.checkIntervalSec) * time.Second)
		}
	}
}
