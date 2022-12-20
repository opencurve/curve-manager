package baserpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/opencurve/pigeon"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	GBaseClient *BaseRpc
)

const (
	RPC_TIMEOUT_MS  = "rpc.timeout.ms"
	RPC_RETRY_TIMES = "rpc.retry.times"

	DEFAULT_RPC_TIMEOUT_MS  = 500
	DEFAULT_RPC_RETRY_TIMES = 3
)

type BaseRpc struct {
	timeout    time.Duration
	retryTimes uint32
	lock       sync.RWMutex
	conns      map[string]*grpc.ClientConn
}

type RpcContext struct {
	addrs []string // endpoint: 127.0.0.1:6666
	name  string
}

func NewRpcContext(addrs []string, funcName string) *RpcContext {
	return &RpcContext{
		addrs: addrs,
		name:  funcName,
	}
}

type Rpc interface {
	NewRpcClient(cc grpc.ClientConnInterface)
	Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error)
}

type RpcResult struct {
	Addr   string
	Err    error
	Result interface{}
}

func Init(cfg *pigeon.Configure) {
	timeout := cfg.GetConfig().GetInt(RPC_TIMEOUT_MS)
	if timeout == 0 {
		timeout = DEFAULT_RPC_TIMEOUT_MS
	}

	retry := cfg.GetConfig().GetInt(RPC_RETRY_TIMES)
	if retry == 0 {
		retry = DEFAULT_RPC_RETRY_TIMES
	}
	GBaseClient = &BaseRpc{
		timeout:    time.Duration(timeout * int(time.Millisecond)),
		retryTimes: uint32(retry),
		lock:       sync.RWMutex{},
		conns:      make(map[string]*grpc.ClientConn),
	}
}

func (cli *BaseRpc) getOrCreateConn(addr string, ctx context.Context) (*grpc.ClientConn, error) {
	cli.lock.RLock()
	conn, ok := cli.conns[addr]
	cli.lock.RUnlock()
	if ok {
		return conn, nil
	}

	cli.lock.Lock()
	defer cli.lock.Unlock()
	conn, ok = cli.conns[addr]
	if ok {
		return conn, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), cli.timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	cli.conns[addr] = conn
	return conn, nil
}

func (cli *BaseRpc) SendRpc(ctx *RpcContext, rpcFunc Rpc) *RpcResult {
	size := len(ctx.addrs)
	if size == 0 {
		return &RpcResult{
			Addr:   "",
			Err:    fmt.Errorf("empty addr"),
			Result: nil,
		}
	}
	results := make(chan RpcResult, size)
	for _, addr := range ctx.addrs {
		go func(address string) {
			ctx, cancel := context.WithTimeout(context.Background(), cli.timeout)
			defer cancel()
			conn, err := cli.getOrCreateConn(address, ctx)
			if err != nil {
				results <- RpcResult{
					Addr:   address,
					Err:    err,
					Result: nil,
				}
			} else {
				rpcFunc.NewRpcClient(conn)
				res, err := rpcFunc.Stub_Func(ctx, grpc_retry.WithMax(uint(cli.retryTimes)),
					grpc_retry.WithCodes(codes.Unknown, codes.Unavailable, codes.DeadlineExceeded))
				results <- RpcResult{
					Addr:   address,
					Err:    err,
					Result: res,
				}
			}
		}(addr)
	}
	count := 0
	var rpcErr string
	for res := range results {
		if res.Err == nil {
			return &res
		}
		count = count + 1
		rpcErr = fmt.Sprintf("%s;%s:%s", rpcErr, res.Addr, res.Err.Error())
		if count >= size {
			break
		}
	}
	return &RpcResult{
		Addr:   "",
		Err:    fmt.Errorf(rpcErr),
		Result: nil,
	}
}
