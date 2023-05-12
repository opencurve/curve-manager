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
* Created Date: 2023-04-13
* Author: wanghai (SeanHai)
 */

package deploy

import (
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/pigeon"
)

const (
	MODULE = "deploy"
)

var METHOD_REQUEST map[string]Request

type (
	HandlerFunc func(r *pigeon.Request, ctx *Context) bool

	Context struct {
		Data interface{}
	}

	Request struct {
		httpMethod string
		method     string
		vType      interface{}
		handler    HandlerFunc
	}
)

func init() {
	METHOD_REQUEST = map[string]Request{}
	for _, request := range requests {
		METHOD_REQUEST[request.method] = request
		core.BELONG[request.method] = MODULE
	}
}

type ListHostRequest struct{}

type PreCommitHostRequest struct {
	Hosts string `json:"hosts" binding:"required"`
}

type CommitHostRequest struct {
	Hosts string `json:"hosts" binding:"required"`
}

type ListDiskRequest struct{}

type CommitDiskRequest struct {
	Disks string `json:"disks" binding:"required"`
}

type GetFormatStatusRequest struct{}

type FormatDiskRequest struct{}

type ShowConfigRequest struct{}

type CommitConfigRequest struct {
	Name string `json:"name" binding:"required"`
	Conf string `json:"conf" binding:"required"`
}

type ListClusterRequest struct{}

type CheckoutClusterRequest struct {
	Name string `json:"name" binding:"required"`
}

type AddClusterRequest struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc"`
	Topo string `json:"topo"`
}

type DeployClusterRequest struct{}

var requests = []Request{
	{
		core.HTTP_GET,
		core.DEPLOY_HOST_LIST,
		ListHostRequest{},
		DealDeploy,
	},
	{
		core.HTTP_POST,
		core.DEPLOY_HOST_COMMIT,
		CommitHostRequest{},
		DealDeploy,
	},
	{
		core.HTTP_GET,
		core.DEPLOY_DISK_LIST,
		ListDiskRequest{},
		DealDeploy,
	},
	{
		core.HTTP_POST,
		core.DEPLOY_DISK_COMMIT,
		CommitDiskRequest{},
		DealDeploy,
	},
	{
		core.HTTP_GET,
		core.DEPLOY_DISK_FORMAT_STATUS,
		GetFormatStatusRequest{},
		DealDeploy,
	},
	{
		core.HTTP_GET,
		core.DEPLOY_DISK_FORMAT,
		FormatDiskRequest{},
		DealDeploy,
	},
	{
		core.HTTP_GET,
		core.DEPLOY_CONFIG_SHOW,
		ShowConfigRequest{},
		DealDeploy,
	},
	{
		core.HTTP_POST,
		core.DEPLOY_CONFIG_COMMIT,
		CommitConfigRequest{},
		DealDeploy,
	},
	{
		core.HTTP_GET,
		core.DEPLOY_CLUSTER_LIST,
		ListClusterRequest{},
		DealDeploy,
	},
	{
		core.HTTP_POST,
		core.DEPLOY_CLUSTER_CHECKOUT,
		CheckoutClusterRequest{},
		DealDeploy,
	},
	{
		core.HTTP_POST,
		core.DEPLOY_CLUSTER_ADD,
		AddClusterRequest{},
		DealDeploy,
	},
	{
		core.HTTP_GET,
		core.DEPLOY_CLUSTER_DEPLOY,
		DeployClusterRequest{},
		DealDeploy,
	},
}
