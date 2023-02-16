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

package baserpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/shimingyah/pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	GBaseClient *BaseRpc
)

const (
	DEFAULT_RPC_TIMEOUT_MS  = 500
	DEFAULT_RPC_RETRY_TIMES = 3
)

type RpcResult common.QueryResult

type BaseRpc struct {
	timeout    time.Duration
	retryTimes uint32
	lock       sync.RWMutex
	connPools  map[string]pool.Pool
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

func init() {
	GBaseClient = &BaseRpc{
		timeout:    time.Duration(DEFAULT_RPC_TIMEOUT_MS * int(time.Millisecond)),
		retryTimes: uint32(DEFAULT_RPC_RETRY_TIMES),
		lock:       sync.RWMutex{},
		connPools:  make(map[string]pool.Pool),
	}
}

func (cli *BaseRpc) getOrCreateConn(addr string, ctx context.Context) (*grpc.ClientConn, error) {
	// cli.lock.RLock()
	// cpool, ok := cli.connPools[addr]
	// cli.lock.RUnlock()
	// if ok {
	// 	conn, err := cpool.Get()
	// 	if err != nil {
	// 		return nil, fmt.Errorf("conn pool get conn failed, addr: %s, err: %v", addr, err)
	// 	}
	// 	return conn.Value(), nil
	// }

	// cli.lock.Lock()
	// defer cli.lock.Unlock()
	// cpool, ok = cli.connPools[addr]
	// if ok {
	// 	conn, err := cpool.Get()
	// 	if err != nil {
	// 		return nil, fmt.Errorf("conn pool get conn failed, addr: %s, err: %v", addr, err)
	// 	}
	// 	return conn.Value(), nil
	// }

	// p, err := pool.New(addr, pool.DefaultOptions)
	// if err != nil {
	// 	return nil, fmt.Errorf("new conn pool failed, addr: %s, err: %v", addr, err)
	// }

	// cli.connPools[addr] = p
	// conn, err := p.Get()
	// if err != nil {
	// 	return nil, fmt.Errorf("conn pool get conn failed, addr: %s, err: %v", addr, err)
	// }
	// return conn.Value(), nil
	ctx, cancel := context.WithTimeout(context.Background(), cli.timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (cli *BaseRpc) SendRpc(ctx *RpcContext, rpcFunc Rpc) *RpcResult {
	size := len(ctx.addrs)
	if size == 0 {
		return &RpcResult{
			Key:    "",
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
					Key:    address,
					Err:    err,
					Result: nil,
				}
			} else {
				rpcFunc.NewRpcClient(conn)
				res, err := rpcFunc.Stub_Func(ctx, grpc_retry.WithMax(uint(cli.retryTimes)),
					grpc_retry.WithCodes(codes.Unknown, codes.Unavailable, codes.DeadlineExceeded))
				results <- RpcResult{
					Key:    address,
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
		count += 1
		rpcErr = fmt.Sprintf("%s;%s:%s", rpcErr, res.Key, res.Err.Error())
		if count >= size {
			break
		}
	}
	return &RpcResult{
		Key:    "",
		Err:    fmt.Errorf(rpcErr),
		Result: nil,
	}
}
