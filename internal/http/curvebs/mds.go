package curvebs

import (
	"github.com/opencurve/curve-manager/internal/http/baseHttp"
	"time"
)

type MdsClientOption struct {
	TimeoutMs  int
	RetryTimes uint32
	Addrs      []string
}

type MdsClient struct {
	addrs           []string
	baseClient_http baseHttp.BaseHttp
}

func NewMdsClient(option MdsClientOption) *MdsClient {
	return &MdsClient{
		addrs: option.Addrs,
		baseClient_http: baseHttp.BaseHttp{
			Timeout:    time.Duration(option.TimeoutMs * int(time.Millisecond)),
			RetryTimes: option.RetryTimes,
		},
	}
}
