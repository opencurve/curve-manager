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

	// 401
	USER_IS_UNAUTHORIZED = Errno{401001, "user is unauthorized"}

	// 403
	REQUEST_IS_DENIED_FOR_SIGNATURE = Errno{403000, "request is denied for signature"}
	GET_USER_FAILED                 = Errno{403001, "user not exist"}
	USER_PASSWORD_NOT_MATCH         = Errno{403002, "user password not match"}
	WRITE_USER_LOGIN_FAILED         = Errno{403003, "only permit one write user login"}
	OPERATION_IS_NOT_PERMIT         = Errno{403004, "operation is not permitted"}

	// 405
	UNSUPPORT_HTTP_METHOD = Errno{405001, "unsupport http method"}

	// 503
	UNKNOW_ERROR = Errno{503001, "unknown error"}

	// user/storage
	CREATE_USER_FAILED            = Errno{503002, "create user failed"}
	DELETE_USER_FAILED            = Errno{503003, "delete user failed"}
	GET_USER_PASSWORD_FAILED      = Errno{503004, "get user password failed"}
	UPDATE_USER_PASSWORD_FAILED   = Errno{503005, "update user password failed"}
	GET_USER_EMAIL_FAILED         = Errno{503006, "get user email failed"}
	USER_EMAIL_EMPTY              = Errno{503007, "user email is empty"}
	SEND_USER_PASSWORD_FAILED     = Errno{503008, "send user email failed"}
	LIST_USER_FAILED              = Errno{503009, "list user failed"}
	UPDATE_USER_EMAIL_FAILED      = Errno{503010, "update user email failed"}
	UPDATE_USER_PERMISSION_FAILED = Errno{503011, "update user permission failed"}

	// hadware/metric
	GET_INSTANCE_BY_HOSTNAME_FAILED  = Errno{503101, "get instance by hostname failed"}
	GET_HOSTNAME_BY_INSTANCE_FAILED  = Errno{503102, "get hostname by instance failed"}
	LIST_DISK_INFO_FAILED            = Errno{503103, "list disk info failed"}
	GET_DISK_FILESYSTEM_INFO_FAILED  = Errno{503104, "get disk filesystem info failed"}
	GET_DISK_TYPE_FAILED             = Errno{503105, "get disk type failed"}
	GET_DISK_WRITE_CACHE_FAILED      = Errno{503106, "get disk write cache failed"}
	GET_HOST_INFO_FAILED             = Errno{503107, "get host info failed"}
	GET_HOST_CPU_INFO_FAILED         = Errno{503108, "get host cpu info failed"}
	GET_HOST_MEM_INFO_FAILED         = Errno{503109, "get host memory info failed"}
	GET_HOST_CPU_UTILIZATION_FAILED  = Errno{503110, "get host cpu utilization failed"}
	GET_HOST_MEM_UTILIZATION_FAILED  = Errno{503111, "get host momery utilization failed"}
	GET_HOST_DISK_PERFORMANCE_FAILED = Errno{503112, "get host disk performance failed"}
	GET_HOST_NETWORK_TRAFFIC_FAILED  = Errno{503113, "get host network traffic failed"}

	// curve/metric
	GET_POOL_ITEM_NUMBER_FAILED    = Errno{503201, "get pool item number failed"}
	GET_POOL_PERFORMANCE_FAILED    = Errno{503202, "get pool performance failed"}
	GET_ROOT_AUTH_FAILED           = Errno{503203, "get root auth failed"}
	GET_VOLUME_PERFORMANCE_FAILED  = Errno{503204, "get volume performance failed"}
	GET_CLUSTER_PERFORMANCE_FAILED = Errno{503205, "get cluster performance failed"}
	GET_CLUSTER_SPACE_FAILED       = Errno{503206, "get cluster space failed"}

	GET_CHUNKSERVER_IN_CLUSTER_FAILED = Errno{503301, "get chunkserver in cluster failed"}
	LIST_POOL_FAILED                  = Errno{503302, "list pool failed"}
	GET_POOL_FAILED                   = Errno{503303, "get pool failed"}
	GET_POOL_SPACE_FAILED             = Errno{503304, "get pool space failed"}
	LIST_POOL_ZONE_FAILED             = Errno{503305, "list pool zone failed"}
	LIST_VOLUME_FAILED                = Errno{503306, "list volume failed"}
	GET_VOLUME_INFO_FAILED            = Errno{503307, "get volume failed"}
	GET_VOLUME_ALLOC_SIZE_FAILED      = Errno{503308, "get volume alloc size failed"}
	GET_VOLUME_SIZE_FAILED            = Errno{503309, "get volume size failed"}
	DELETE_VOLUME_FAILED              = Errno{503310, "delete volume failed"}
	CREATE_VOLUME_FAILED              = Errno{503311, "create volume failed"}
	EXTEND_VOLUME_FAILED              = Errno{503312, "extend volume failed"}
	VOLUME_THROTTLE_FAILED            = Errno{503313, "update volume throttle params failed"}
	RECOVER_VOLUME_FAILED             = Errno{503314, "recover volume failed"}
	CLONE_VOLUME_FAILED               = Errno{503315, "create clone volume failed"}
	CREATE_SNAPSHOT_FAILED            = Errno{503316, "create snapshot failed"}
	LIST_SNAPSHOT_FAILED              = Errno{503317, "list snapshot failed"}
	CANCEL_SNAPSHOT_FAILED            = Errno{503318, "cancel snapshot failed"}
	GET_CLONE_TASKS_FAILED            = Errno{503319, "get clone tasks failed"}
	FLATTEN_FAILED                    = Errno{503320, "flatten failed"}
	GET_SNAPSHOT_FAILED               = Errno{503321, "get snapshot failed"}
	DELETE_SNAPSHOT_FAILED            = Errno{503322, "delete snapshot failed"}
	FIND_VOLUME_MOUNTPOINT_FAILED     = Errno{503323, "find volume mountpoint failed"}
	GET_SYSTEM_LOG_FAILED             = Errno{503324, "get system operation log failed"}
	GET_SYSTEM_ALERT_FAILED           = Errno{503325, "get system alert failed"}
	GET_UNREAD_ALERT_NUM_FAILED       = Errno{503326, "get unread alert num failed"}
	GET_READ_ALERT_ID_FAILED          = Errno{503327, "get read alert id failed"}
	ADD_READ_ALERT_ID_FAILED          = Errno{503328, "add read alert id failed"}
	UPDATE_UNREAD_ALERT_ID_FAILED     = Errno{503329, "update unread alert id failed"}
	GET_ALERT_CONF_FAILED             = Errno{503330, "get alert conf failed"}
	UPDATE_ALERT_CONF_FAILED          = Errno{503331, "update alert conf failed"}
	GET_ALERT_USER_FAILED             = Errno{503332, "get alert user failed"}
	ADD_ALERT_USER_FAILED             = Errno{503333, "add alert user failed"}
	DELETE_ALERT_USER_FAILED          = Errno{503334, "delete alert user failed"}
	LIST_USER_WITH_EMAIL_FAILED       = Errno{503335, "list user with email failed"}
)
