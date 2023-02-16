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
* Created Date: 2023-02-11
* Author: wanghai (SeanHai)
 */

package errno

type IErrno interface {
	Code() int
	HTTPCode() int
	Description() string
}

type Errno struct {
	code        int
	description string
}

func (e Errno) Code() int           { return e.code }
func (e Errno) HTTPCode() int       { return e.code / 1000 }
func (e Errno) Description() string { return e.description }

var (
	OK = Errno{0, "success"}

	// 400
	UNSUPPORT_REQUEST_URI     = Errno{400001, "unsupport request uri"}
	UNSUPPORT_METHOD_ARGUMENT = Errno{400002, "unsupport method argument"}
	HTTP_METHOD_MISMATCHED    = Errno{400003, "http method mismatch"}
	BAD_REQUEST_FORM_PARAM    = Errno{400004, "bad request form param"}

	// 403
	REQUEST_IS_DENIED_FOR_SIGNATURE = Errno{403000, "request is denied for signature"}

	// 405
	UNSUPPORT_HTTP_METHOD = Errno{405001, "unsupport http method"}

	// 503
	// user/storage
	GET_USER_FAILED             = Errno{503001, "get user failed"}
	USER_PASSWORD_NOT_MATCH     = Errno{503002, "user password not match"}
	CREATE_USER_FAILED          = Errno{503003, "create user failed"}
	DELETE_USER_FAILED          = Errno{503004, "delete user failed"}
	GET_USER_PASSWORD_FAILED    = Errno{503005, "get user password failed"}
	UPDATE_USER_PASSWORD_FAILED = Errno{503006, "update user password failed"}
	GET_USER_EMAIL_FAILED       = Errno{503007, "get user email failed"}
	USER_EMAIL_EMPTY            = Errno{503008, "user email is empty"}
	SEND_USER_PASSWORD_FAILED   = Errno{503009, "send user email failed"}
	UPDATE_USER_INFO_FAILED     = Errno{503010, "update user info failed"}
	LIST_USER_FAILED            = Errno{503011, "list user failed"}

	// hadware/metric
	GET_INSTANCE_BY_HOSTNAME_FAILED  = Errno{503101, "get instance by hostname failed"}
	GET_HOSTNAME_BY_INSTANCE_FAILED  = Errno{503102, "get hostname by instance failed"}
	LIST_DISK_FAILED                 = Errno{503103, "list disk failed"}
	GET_FILESYSTEM_INFO_FAILED       = Errno{503104, "get filesystem info failed"}
	GET_HOST_INFO_FAILED             = Errno{503105, "get host info failed"}
	GET_HOST_CPU_INFO_FAILED         = Errno{503106, "get host cpu info failed"}
	GET_HOST_MEM_INFO_FAILED         = Errno{503107, "get host memory info failed"}
	GET_HOST_DISK_NUM_FAILED         = Errno{503108, "get host disk number failed"}
	GET_HOST_CPU_UTILIZATION_FAILED  = Errno{503109, "get host cpu utilization failed"}
	GET_HOST_MEM_UTILIZATION_FAILED  = Errno{503110, "get host momery utilization failed"}
	GET_HOST_DISK_PERFORMANCE_FAILED = Errno{503111, "get host disk performance failed"}
	GET_HOST_NETWORK_TRAFFIC_FAILED  = Errno{503112, "get host network traffic failed"}

	// curve/metric
	GET_ETCD_STATUS_FAILED           = Errno{503201, "get etcd status failed"}
	GET_MDS_STATUS_FAILED            = Errno{503202, "get mds status failed"}
	GET_SNAPSHOT_CLONE_STATUS_FAILED = Errno{503203, "get snapshotcloneserver status failed"}
	GET_CHUNKSERVER_VERSION_FAILED   = Errno{503204, "get host network traffic failed"}
	GET_POOL_ITEM_NUMBER_FAILED      = Errno{503205, "get pool item number failed"}
	GET_POOL_PERFORMANCE_FAILED      = Errno{503206, "get pool performance failed"}
	GET_ROOT_AUTH_FAILED             = Errno{503207, "get root auth failed"}
	GET_VOLUME_PERFORMANCE_FAILED    = Errno{503208, "get volume performance failed"}

	// rpc
	LIST_SNAPSHOT_FAILED              = Errno{503301, "list snapshot failed"}
	GET_CHUNKSERVER_IN_CLUSTER_FAILED = Errno{503302, "get chunkserver in cluster failed"}
	LIST_POOL_FAILED                  = Errno{503303, "list pool failed"}
	GET_POOL_SPACE_FAILED             = Errno{503304, "get pool space failed"}
	GET_CLUSTER_PERFORMANCE_FAILED    = Errno{503305, "get cluster performance failed"}
	LIST_POOL_ZONE_FAILED             = Errno{503306, "list pool zone failed"}
	LIST_VOLUME_FAILED                = Errno{503307, "list volume failed"}
	GET_VOLUME_ALLOC_SIZE_FAILED      = Errno{503308, "get volume alloc size failed"}
)
