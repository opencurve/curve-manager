package curvebs

import (
	"github.com/opencurve/curve-manager/internal/common"
	"strings"
)

var (
	GMdsClient *MdsClient
)

const (
	CURVEBS_MDS_ADDRESS = "mds.address"

	DEFAULT_RPC_TIMEOUT_MS  = 500
	DEFAULT_RPC_RETRY_TIMES = 3
)

func Init(cfg map[string]string) {
	addrs := cfg[CURVEBS_MDS_ADDRESS]
	GMdsClient = NewMdsClient(MdsClientOption{
		TimeoutMs:  DEFAULT_RPC_TIMEOUT_MS,
		RetryTimes: DEFAULT_RPC_RETRY_TIMES,
		Addrs:      strings.Split(addrs, common.CURVEBS_ADDRESS_DELIMITER),
	})
}
