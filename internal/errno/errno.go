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
	USER_LOGIN_FAILED       = Errno{503001, "user login failed"}
	CREATE_USER_FAILED      = Errno{503002, "create user failed"}
	DELETE_USER_FAILED      = Errno{503003, "delete user failed"}
	CHANGE_PASSWORD_FAILED  = Errno{503004, "change user password failed"}
	UPDATE_USER_INFO_FAILED = Errno{503005, "update user info failed"}
	LIST_USER_FAILED        = Errno{503006, "list user failed"}

	GET_CLUSTER_STATUS_FAILED        = Errno{503101, "get cluster status failed"}
	GET_ETCD_STATUS_FAILED           = Errno{503102, "get etcd status failed"}
	GET_MDS_STATUS_FAILED            = Errno{503103, "get mds status failed"}
	GET_SNAPSHOT_CLONE_STATUS_FAILED = Errno{503104, "get snapshotcloneserver status failed"}
	GET_CHUNKSERVER_STATUS_FAILED    = Errno{503105, "get chunkserver status failed"}

	GET_CLUSTER_SPACE_FAILED       = Errno{503201, "get cluster space failed"}
	GET_CLUSTER_PERFORMANCE_FAILED = Errno{503202, "get cluster performance failed"}
	LIST_TOPO_FAILED               = Errno{503203, "list topo failed"}
	LIST_POOL_FAILED               = Errno{503204, "list pool failed"}
	LIST_VOLUME_FAILED             = Errno{503205, "list volume failed"}
)
