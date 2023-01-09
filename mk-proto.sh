# curvebs proto

## common.proto
protoc --go_out=internal --proto_path=external/_curve \
    external/_curve/proto/common.proto

## chunk.proto
protoc --go_out=internal --proto_path=external/_curve \
    external/_curve/proto/chunk.proto

## chunkserver.proto
protoc --go_out=internal --proto_path=external/_curve \
    external/_curve/proto/chunkserver.proto

## cli.proto
protoc --go_out=internal --proto_path=external/_curve \
    external/_curve/proto/cli.proto

## cli2.proto
protoc --go_out=internal --proto_path=external/_curve \
    --go_opt=Mproto/common.proto=github.com/opencurve/curve-manager/internal/proto/common \
    external/_curve/proto/cli2.proto

## configuration.proto
protoc --go_out=internal --proto_path=external/_curve \
    external/_curve/proto/configuration.proto

## copyset.proto
protoc --go_out=internal --proto_path=external/_curve \
    --go_opt=Mproto/common.proto=github.com/opencurve/curve-manager/internal/proto/common \
    external/_curve/proto/copyset.proto

## curve_storage.proto
protoc --go_out=internal --proto_path=external/_curve \
    external/_curve/proto/curve_storage.proto

## scan.proto
protoc --go_out=internal --proto_path=external/_curve \
    external/_curve/proto/scan.proto

## heartbeat.proto
protoc --go_out=internal --proto_path=external/_curve \
    --go_opt=Mproto/common.proto=github.com/opencurve/curve-manager/internal/proto/common \
    --go_opt=Mproto/scan.proto=github.com/opencurve/curve-manager/internal/proto/scan \
    external/_curve/proto/heartbeat.proto

## integrity.proto
protoc --go_out=internal --proto_path=external/_curve \
    external/_curve/proto/integrity.proto

## namespace2.proto
protoc --go_out=internal --proto_path=external/_curve \
    --go_opt=Mproto/common.proto=github.com/opencurve/curve-manager/internal/proto/common \
    external/_curve/proto/nameserver2.proto

## schedule.proto
protoc --go_out=internal --proto_path=external/_curve \
    external/_curve/proto/schedule.proto

## snapshotcloneserver.proto
protoc --go_out=internal --proto_path=external/_curve \
    external/_curve/proto/snapshotcloneserver.proto

## topology.proto
protoc --go_out=internal --proto_path=external/_curve \
    --go_opt=Mproto/common.proto=github.com/opencurve/curve-manager/internal/proto/common \
    external/_curve/proto/topology.proto

## statuscode.proto
protoc --go_out=internal --proto_path=external/_curve/tools-v2 \
    external/_curve/tools-v2/internal/proto/curvebs/topology/statuscode.proto

protoc --go-grpc_out=internal --proto_path=external/_curve \
    external/_curve/proto/*.proto