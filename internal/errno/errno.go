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
	CREATE_USER_FAILED = Errno{503001, "create user failed"}

	GET_ETCD_STATUS_FAILED = Errno{503101, "get etcd status failed"}

	LIST_PHYSICAL_POOL_FAILED = Errno{503201, "list physical pool failed"}
)
