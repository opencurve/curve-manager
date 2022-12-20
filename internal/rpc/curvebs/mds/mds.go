package mds

import (
	"context"
	"fmt"
	"strings"

	"github.com/opencurve/curve-manager/internal/proto/topology"
	"github.com/opencurve/curve-manager/internal/proto/topology/statuscode"
	"github.com/opencurve/curve-manager/internal/rpc/baserpc"
	comm "github.com/opencurve/curve-manager/internal/rpc/curvebs/common"
	"github.com/opencurve/pigeon"
	"google.golang.org/grpc"
)

var (
	GMdsClient *mdsClient
)

const (
	CURVEBS_MDS_ADDRESS           = "mds.address"
	CURVEBS_MDS_ADDRESS_DELIMITER = ","

	LIST_PHYSICAL_POOL_NAME = "ListPhysicalPool"
)

type mdsClient struct {
	addrs []string
}

func Init(cfg *pigeon.Configure) error {
	addrs := cfg.GetConfig().GetString(CURVEBS_MDS_ADDRESS)
	if len(addrs) == 0 {
		return fmt.Errorf("no cluster mds address found")
	}
	GMdsClient = &mdsClient{
		addrs: strings.Split(addrs, CURVEBS_MDS_ADDRESS_DELIMITER),
	}
	return nil
}

// list physical pool
type ListPhysicalPoolRpc struct {
	ctx            *baserpc.RpcContext
	Request        *topology.ListPhysicalPoolRequest
	topologyClient topology.TopologyServiceClient
}

func (lpRpc *ListPhysicalPoolRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	lpRpc.topologyClient = topology.NewTopologyServiceClient(cc)
}

func (lpRpc *ListPhysicalPoolRpc) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return lpRpc.topologyClient.ListPhysicalPool(ctx, lpRpc.Request, opt...)
}

func (cli *mdsClient) ListPhysicalPools() (interface{}, error) {
	listPoolRpc := &ListPhysicalPoolRpc{}
	listPoolRpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_PHYSICAL_POOL_NAME)
	listPoolRpc.Request = &topology.ListPhysicalPoolRequest{}
	ret := baserpc.GBaseClient.SendRpc(listPoolRpc.ctx, listPoolRpc)
	if ret.Err != nil {
		return nil, ret.Err
	}

	response := ret.Result.(*topology.ListPhysicalPoolResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return nil, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}

	var infos []comm.PhysicalPool
	for _, pool := range response.GetPhysicalPoolInfos() {
		info := comm.PhysicalPool{}
		info.Id = pool.GetPhysicalPoolID()
		info.Name = pool.GetPhysicalPoolName()
		info.Desc = pool.GetDesc()
		infos = append(infos, info)
	}
	return infos, nil
}
