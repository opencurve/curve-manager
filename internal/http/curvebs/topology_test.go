package curvebs

import (
	"testing"
)

var (
	clientOption MdsClientOption = MdsClientOption{
		TimeoutMs:  500,
		RetryTimes: 3,
		Addrs:      []string{"192.168.170.138:6702"},
	}
)

func TestMdsClient_ListPhysicalPool_http(t *testing.T) {
	MdsClient := NewMdsClient(clientOption)
	pools, err := MdsClient.ListPhysicalPool_http()
	print(pools)
	print(err)
}
