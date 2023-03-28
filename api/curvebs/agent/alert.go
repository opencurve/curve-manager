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
	"strconv"
	"sync"
	"time"

	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/storage"
	"github.com/opencurve/pigeon"
)

const (
	ALERT_REQUEST_ID = "alert"

	ALERT_CLUSTER               = "cluster"
	ALERT_SPACE                 = "space"
	ALERT_ETCD                  = "etcd"
	ALERT_MDS                   = "mds"
	ALERT_CHUNKSERVER           = "chunkserver"
	ALERT_SNAPSHOT_CLONE_SERVER = "snapshotcloneserver"

	UPDATE_ALERT_CONF_INTERVAL_SEC = 1 * 30

	CLUSTER_ALERT_INTERVAL_SEC  = 1 * 60
	CLUSTER_ALERT_TRIGGER_TIMES = 3
	CLUSTER_ALERT_RULE          = "cluster is not healthy"
	CLUSTER_ALERT_DESC          = "check cluster healthy"

	CAPACITY_ALERT_INTERVAL_SEC  = 1 * 60 * 60
	CAPACITY_ALERT_TRIGGER_TIMES = 1
	CAPACITY_ALERT_LIMIT_PERCENT = 80
	CAPACITY_ALERT_DESC          = "check cluster space used"

	SERVICE_ALERT_INTERVAL_SEC  = 1 * 60
	SERVICE_ALERT_TRIGGER_TIMES = 1
	SERVICE_ALERT_RULE          = "leaderNum != 1 or offlineNum > 0"
	SERVICE_ALERT_DESC          = "check service status"
)

type AlertConf struct {
	Name       string   `json:"name"`
	Level      string   `json:"level"`
	Interval   uint32   `json:"interval"`
	Times      uint32   `json:"times"`
	Enable     bool     `json:"enable"`
	Rule       string   `json:"rule"`
	Desc       string   `json:"desc"`
	AlertUsers []string `json:"alertUsers"`
}

var (
	alertClients map[string]Alert = make(map[string]Alert)

	defaultAlertConf = []storage.AlertConf{
		{
			Name:       ALERT_CLUSTER,
			Level:      storage.ALERT_WARNING,
			LevelStr:   storage.WARNING,
			Interval:   CLUSTER_ALERT_INTERVAL_SEC,
			Times:      CLUSTER_ALERT_TRIGGER_TIMES,
			Enable:     1,
			EnableBool: true,
			Rule:       CLUSTER_ALERT_RULE,
			Desc:       CLUSTER_ALERT_DESC,
		},
		{
			Name:       ALERT_SPACE,
			Level:      storage.ALERT_WARNING,
			LevelStr:   storage.WARNING,
			Interval:   CAPACITY_ALERT_INTERVAL_SEC,
			Times:      CAPACITY_ALERT_TRIGGER_TIMES,
			Enable:     1,
			EnableBool: true,
			Rule:       strconv.FormatUint(CAPACITY_ALERT_LIMIT_PERCENT, 10),
			Desc:       CAPACITY_ALERT_DESC,
		},
		{
			Name:       ALERT_ETCD,
			Level:      storage.ALERT_CRITICAL,
			LevelStr:   storage.CRITICAL,
			Interval:   SERVICE_ALERT_INTERVAL_SEC,
			Times:      SERVICE_ALERT_TRIGGER_TIMES,
			Enable:     1,
			EnableBool: true,
			Rule:       SERVICE_ALERT_RULE,
			Desc:       SERVICE_ALERT_DESC,
		},
		{
			Name:       ALERT_MDS,
			Level:      storage.ALERT_CRITICAL,
			LevelStr:   storage.CRITICAL,
			Interval:   SERVICE_ALERT_INTERVAL_SEC,
			Times:      SERVICE_ALERT_TRIGGER_TIMES,
			Enable:     1,
			EnableBool: true,
			Rule:       SERVICE_ALERT_RULE,
			Desc:       SERVICE_ALERT_DESC,
		},
		{
			Name:       ALERT_SNAPSHOT_CLONE_SERVER,
			Level:      storage.ALERT_CRITICAL,
			LevelStr:   storage.CRITICAL,
			Interval:   SERVICE_ALERT_INTERVAL_SEC,
			Times:      SERVICE_ALERT_TRIGGER_TIMES,
			Enable:     1,
			EnableBool: true,
			Rule:       SERVICE_ALERT_RULE,
			Desc:       SERVICE_ALERT_DESC,
		},
		{
			Name:       ALERT_CHUNKSERVER,
			Level:      storage.ALERT_CRITICAL,
			LevelStr:   storage.CRITICAL,
			Interval:   SERVICE_ALERT_INTERVAL_SEC,
			Times:      SERVICE_ALERT_TRIGGER_TIMES,
			Enable:     1,
			EnableBool: true,
			Rule:       SERVICE_ALERT_RULE,
			Desc:       SERVICE_ALERT_DESC,
		},
	}
)

type alertContext struct {
	opt   storage.AlertConf
	mutex sync.Mutex
	times uint32
}

func (ctx *alertContext) getEnable() bool {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.opt.EnableBool
}

func (ctx *alertContext) setEnable(enable bool) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.opt.EnableBool = enable
}

func (ctx *alertContext) getInterval() uint32 {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.opt.Interval
}

func (ctx *alertContext) setInterval(interval uint32) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.opt.Interval = interval
}

func (ctx *alertContext) getTimes() uint32 {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.opt.Times
}

func (ctx *alertContext) setTimes(times uint32) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.opt.Times = times
}

func (ctx *alertContext) getRule() string {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.opt.Rule
}

func (ctx *alertContext) setRule(rule string) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.opt.Rule = rule
}

type Alert interface {
	check(logger *pigeon.Logger)
}

type clusterAlert struct {
	ctx alertContext
}

type spaceAlert struct {
	ctx alertContext
}

type etcdServiceAlert struct {
	ctx alertContext
}

type mdsServiceAlert struct {
	ctx alertContext
}

type snapshotCloneServiceAlert struct {
	ctx alertContext
}

type chunkserverServiceAlert struct {
	ctx alertContext
}

func initAlerts(logger *pigeon.Logger) error {
	var initConfs []storage.AlertConf
	alertInfo, err := storage.GetAlertConf()
	if err != nil {
		return err
	}
	if len(alertInfo) == 0 {
		for _, conf := range defaultAlertConf {
			err := storage.AddAlertConf(&conf)
			if err != nil {
				return err
			}
		}
		initConfs = defaultAlertConf
	} else {
		initConfs = alertInfo
	}

	for _, conf := range initConfs {
		switch conf.Name {
		case ALERT_CLUSTER:
			alertClients[ALERT_CLUSTER] = &clusterAlert{
				ctx: alertContext{
					opt:   conf,
					mutex: sync.Mutex{},
					times: 0,
				},
			}
			go alertClients[ALERT_CLUSTER].check(logger)
		case ALERT_SPACE:
			alertClients[ALERT_SPACE] = &spaceAlert{
				ctx: alertContext{
					opt:   conf,
					mutex: sync.Mutex{},
					times: 0,
				},
			}
			go alertClients[ALERT_SPACE].check(logger)
		case ALERT_ETCD:
			alertClients[ALERT_ETCD] = &etcdServiceAlert{
				ctx: alertContext{
					opt:   conf,
					mutex: sync.Mutex{},
					times: 0,
				},
			}
			go alertClients[ALERT_ETCD].check(logger)
		case ALERT_MDS:
			alertClients[ALERT_MDS] = &mdsServiceAlert{
				ctx: alertContext{
					opt:   conf,
					mutex: sync.Mutex{},
					times: 0,
				},
			}
			go alertClients[ALERT_MDS].check(logger)
		case ALERT_CHUNKSERVER:
			alertClients[ALERT_CHUNKSERVER] = &chunkserverServiceAlert{
				ctx: alertContext{
					opt:   conf,
					mutex: sync.Mutex{},
					times: 0,
				},
			}
			go alertClients[ALERT_CHUNKSERVER].check(logger)
		case ALERT_SNAPSHOT_CLONE_SERVER:
			alertClients[ALERT_SNAPSHOT_CLONE_SERVER] = &snapshotCloneServiceAlert{
				ctx: alertContext{
					opt:   conf,
					mutex: sync.Mutex{},
					times: 0,
				},
			}
			go alertClients[ALERT_SNAPSHOT_CLONE_SERVER].check(logger)
		}
	}
	go updateAlertConf(logger)
	return nil
}

func updateAlertConf(logger *pigeon.Logger) {
	timer := time.NewTimer(time.Duration(UPDATE_ALERT_CONF_INTERVAL_SEC) * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			alertInfo, err := storage.GetAlertConf()
			if err != nil {
				logger.Error("UpdateAlertConf get alert conf failed",
					pigeon.Field("error", err),
					pigeon.Field("requestId", ALERT_REQUEST_ID))
			} else {
				for _, conf := range alertInfo {
					if _, ok := alertClients[conf.Name]; ok {
						switch conf.Name {
						case ALERT_CLUSTER:
							updateConf(&alertClients[conf.Name].(*clusterAlert).ctx, &conf)
						case ALERT_SPACE:
							updateConf(&alertClients[conf.Name].(*spaceAlert).ctx, &conf)
						case ALERT_ETCD:
							updateConf(&alertClients[conf.Name].(*etcdServiceAlert).ctx, &conf)
						case ALERT_MDS:
							updateConf(&alertClients[conf.Name].(*mdsServiceAlert).ctx, &conf)
						case ALERT_SNAPSHOT_CLONE_SERVER:
							updateConf(&alertClients[conf.Name].(*snapshotCloneServiceAlert).ctx, &conf)
						case ALERT_CHUNKSERVER:
							updateConf(&alertClients[conf.Name].(*chunkserverServiceAlert).ctx, &conf)
						}
					}
				}
			}
			timer.Reset(time.Duration(UPDATE_ALERT_CONF_INTERVAL_SEC) * time.Second)
		}
	}
}

func updateConf(ctx *alertContext, conf *storage.AlertConf) {
	ctx.setEnable(conf.EnableBool)
	ctx.setInterval(conf.Interval)
	ctx.setRule(conf.Rule)
	ctx.setTimes(conf.Times)
}

func handleAlert(level int, name string, duration uint32, summary string) error {
	now := time.Now().UnixMilli()
	alert := &storage.Alert{
		TimeMs:      now,
		Level:       level,
		Name:        name,
		DurationSec: duration,
		Summary:     summary,
	}
	err := storage.AddAlert(alert)
	if err != nil {
		return err
	}
	err = storage.SendAlert(alert)
	return err
}

func (alert *clusterAlert) check(logger *pigeon.Logger) {
	timer := time.NewTimer(time.Duration(alert.ctx.getInterval()) * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			if alert.ctx.getEnable() {
				status := GetClusterStatus(logger, ALERT_REQUEST_ID)
				if !status.(ClusterStatus).Healthy {
					alert.ctx.times++
					if alert.ctx.times >= alert.ctx.getTimes() {
						summary := "cluster is not healthy"
						unhealthy := status.(ClusterStatus).CopysetNum.Unhealthy
						if unhealthy > 0 {
							summary = fmt.Sprintf("%s, unhealthy copysets number is %d", summary, unhealthy)
						}
						e := handleAlert(alert.ctx.opt.Level, alert.ctx.opt.Name, alert.ctx.getInterval()*alert.ctx.getTimes(), summary)
						if e != nil {
							logger.Error("handle cluster alert info failed",
								pigeon.Field("error", e),
								pigeon.Field("requestId", ALERT_REQUEST_ID))
						}
						alert.ctx.times = 0
					}
				}
			}
			timer.Reset(time.Duration(alert.ctx.getInterval()) * time.Second)
		}
	}
}

func (alert *spaceAlert) check(logger *pigeon.Logger) {
	timer := time.NewTimer(time.Duration(alert.ctx.getInterval()) * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			if alert.ctx.getEnable() {
				space, err := GetClusterSpace(logger, ALERT_REQUEST_ID)
				if err == errno.OK {
					percent := space.(Space).Alloc * 100 / space.(Space).Total
					limitPercent, err := strconv.ParseUint(alert.ctx.getRule(), 10, 64)
					if err != nil {
						logger.Error("space alert parge rule to limit failed, and will use default limit 80%",
							pigeon.Field("rule", alert.ctx.getRule()),
							pigeon.Field("error", err))
						limitPercent = CAPACITY_ALERT_LIMIT_PERCENT
					}
					if percent >= limitPercent {
						alert.ctx.times++
						if alert.ctx.times >= alert.ctx.getTimes() {
							e := handleAlert(alert.ctx.opt.Level, alert.ctx.opt.Name, alert.ctx.getInterval()*alert.ctx.getTimes(),
								fmt.Sprintf("Cluster space have alloced %d%%(alert trigger is %d%%)", percent, limitPercent))
							if e != nil {
								logger.Error("handle space alert falied",
									pigeon.Field("error", e))
							}
							alert.ctx.times = 0
						}
					}
				}
			}
			timer.Reset(time.Duration(alert.ctx.getInterval()) * time.Second)
		}
	}
}

func (alert *etcdServiceAlert) check(logger *pigeon.Logger) {
	timer := time.NewTimer(time.Duration(alert.ctx.getInterval()) * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			if alert.ctx.getEnable() {
				result := checkServiceHealthy(alert.ctx.opt.Name)
				if !result.Result.(serviceStatus).healthy {
					alert.ctx.times++
					if alert.ctx.times >= alert.ctx.getTimes() {
						e := handleAlert(alert.ctx.opt.Level, alert.ctx.opt.Name, alert.ctx.getInterval()*alert.ctx.getTimes(),
							result.Result.(serviceStatus).detail)
						if e != nil {
							logger.Error("handle service alert info failed",
								pigeon.Field("error", e),
								pigeon.Field("requestId", ALERT_REQUEST_ID))
						}
						alert.ctx.times = 0
					}
				}
			}
			timer.Reset(time.Duration(alert.ctx.getInterval()) * time.Second)
		}
	}
}

func (alert *mdsServiceAlert) check(logger *pigeon.Logger) {
	timer := time.NewTimer(time.Duration(alert.ctx.getInterval()) * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			if alert.ctx.getEnable() {
				result := checkServiceHealthy(alert.ctx.opt.Name)
				if !result.Result.(serviceStatus).healthy {
					alert.ctx.times++
					if alert.ctx.times >= alert.ctx.getTimes() {
						e := handleAlert(alert.ctx.opt.Level, alert.ctx.opt.Name, alert.ctx.getInterval()*alert.ctx.getTimes(),
							result.Result.(serviceStatus).detail)
						if e != nil {
							logger.Error("handle service alert info failed",
								pigeon.Field("error", e),
								pigeon.Field("requestId", ALERT_REQUEST_ID))
						}
						alert.ctx.times = 0
					}
				}
			}
			timer.Reset(time.Duration(alert.ctx.getInterval()) * time.Second)
		}
	}
}

func (alert *snapshotCloneServiceAlert) check(logger *pigeon.Logger) {
	timer := time.NewTimer(time.Duration(alert.ctx.getInterval()) * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			if alert.ctx.getEnable() {
				result := checkServiceHealthy(alert.ctx.opt.Name)
				if !result.Result.(serviceStatus).healthy {
					alert.ctx.times++
					if alert.ctx.times >= alert.ctx.getTimes() {
						e := handleAlert(alert.ctx.opt.Level, alert.ctx.opt.Name, alert.ctx.getInterval()*alert.ctx.getTimes(),
							result.Result.(serviceStatus).detail)
						if e != nil {
							logger.Error("handle service alert info failed",
								pigeon.Field("error", e),
								pigeon.Field("requestId", ALERT_REQUEST_ID))
						}
						alert.ctx.times = 0
					}
				}
			}
			timer.Reset(time.Duration(alert.ctx.getInterval()) * time.Second)
		}
	}
}

func (alert *chunkserverServiceAlert) check(logger *pigeon.Logger) {
	timer := time.NewTimer(time.Duration(alert.ctx.getInterval()) * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			if alert.ctx.getEnable() {
				csStatus, err := GetChunkServerStatus(logger, ALERT_REQUEST_ID)
				if err == errno.OK {
					notOnline := csStatus.(*ChunkServerStatus).NotOnlines
					if len(notOnline) > 0 {
						alert.ctx.times++
						if alert.ctx.times >= alert.ctx.getTimes() {
							e := handleAlert(alert.ctx.opt.Level, alert.ctx.opt.Name, alert.ctx.getInterval()*alert.ctx.getTimes(),
								fmt.Sprintf("chunkserver not online number = %d, notOnlines: %v", len(notOnline), notOnline))
							if e != nil {
								logger.Error("handle service alert info failed",
									pigeon.Field("error", e),
									pigeon.Field("requestId", ALERT_REQUEST_ID))
							}
							alert.ctx.times = 0
						}
					}
				}
			}
			timer.Reset(time.Duration(alert.ctx.getInterval()) * time.Second)
		}
	}
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
	return maxId - readId, errno.OK
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

func GetSysAlert(r *pigeon.Request, start, end int64, page, size uint32, name, level, content string) (interface{}, errno.Errno) {
	if start == 0 && end == 0 {
		end = time.Now().UnixMilli()
	}
	info, err := storage.GetAlert(start, end, size, (page-1)*size, name, level, content)
	if err != nil {
		r.Logger().Error("GetAlert failed",
			pigeon.Field("start", start),
			pigeon.Field("end", end),
			pigeon.Field("name", name),
			pigeon.Field("level", level),
			pigeon.Field("content", content),
			pigeon.Field("page", page),
			pigeon.Field("size", size),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_SYSTEM_ALERT_FAILED
	}
	return info, errno.OK
}

func GetAlertConf(r *pigeon.Request) (interface{}, errno.Errno) {
	listConfs := []AlertConf{}
	confs, err := storage.GetAlertConf()
	if err != nil {
		r.Logger().Error("GetAlertConf failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_ALERT_CONF_FAILED
	}
	for _, conf := range confs {
		users, err := storage.GetAlertUser(conf.Name)
		if err != nil {
			r.Logger().Error("GetAlertUser failed",
				pigeon.Field("alertName", conf.Name),
				pigeon.Field("error", err),
				pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
			return nil, errno.GET_ALERT_USER_FAILED
		}
		listConfs = append(listConfs, AlertConf{
			Name:       conf.Name,
			Level:      conf.LevelStr,
			Interval:   conf.Interval,
			Times:      conf.Times,
			Enable:     conf.EnableBool,
			Rule:       conf.Rule,
			Desc:       conf.Desc,
			AlertUsers: users,
		})
	}
	return listConfs, errno.OK
}

func UpdateAlertConf(r *pigeon.Request, enable bool, interval, times uint32, rule string, name string) errno.Errno {
	flag := 1
	if !enable {
		flag = 0
	}
	err := storage.UpdateAlertConf(&storage.AlertConf{
		Interval: interval,
		Enable:   flag,
		Times:    times,
		Rule:     rule,
		Name:     name,
	})
	if err != nil {
		r.Logger().Error("UpdateAlertConf failed",
			pigeon.Field("name", name),
			pigeon.Field("enable", enable),
			pigeon.Field("interval", interval),
			pigeon.Field("times", times),
			pigeon.Field("rule", rule),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return errno.UPDATE_ALERT_CONF_FAILED
	}
	return errno.OK
}

func GetAlertCandidate(r *pigeon.Request) (interface{}, errno.Errno) {
	users, err := storage.ListUserWithEmail()
	if err != nil {
		r.Logger().Error("GetAlertCandidate failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.LIST_USER_WITH_EMAIL_FAILED
	}
	return users, errno.OK
}

func UpdateAlertUser(r *pigeon.Request, alert, user string, op int) errno.Errno {
	if op == 1 {
		err := storage.AddAlertUser(alert, []string{user})
		if err != nil {
			r.Logger().Error("AddAlertUser failed",
				pigeon.Field("alert", alert),
				pigeon.Field("user", user),
				pigeon.Field("error", err),
				pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
			return errno.ADD_ALERT_USER_FAILED
		}
	} else if op == -1 {
		err := storage.DeleteAlertUser(alert, []string{user})
		if err != nil {
			r.Logger().Error("DeleteAlertUser failed",
				pigeon.Field("alert", alert),
				pigeon.Field("user", user),
				pigeon.Field("error", err),
				pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
			return errno.DELETE_ALERT_USER_FAILED
		}
	}
	return errno.OK
}
