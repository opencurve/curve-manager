# curvebs proto

## common.proto
protoc --go_out=internal --proto_path=external/curve \
    external/curve/proto/common.proto

## chunk.proto
protoc --go_out=internal --proto_path=external/curve \
    external/curve/proto/chunk.proto

## chunkserver.proto
protoc --go_out=internal --proto_path=external/curve \
    external/curve/proto/chunkserver.proto

## cli.proto
protoc --go_out=internal --proto_path=external/curve \
    external/curve/proto/chunkserver.proto

## cli2.proto
protoc --go_out=internal --proto_path=external/curve \
    --go_opt=Mproto/common.proto=github.com/opencurve/curve-manager/internal/proto/common \
    external/curve/proto/cli2.proto

## configuration.proto
protoc --go_out=internal --proto_path=external/curve \
    external/curve/proto/configuration.proto

## copyset.proto
protoc --go_out=internal --proto_path=external/curve \
    external/curve/proto/copyset.proto

## curve_storage.proto
protoc --go_out=internal --proto_path=external/curve \
    external/curve/proto/curve_storage.proto

## scan.proto
protoc --go_out=internal --proto_path=external/curve \
    external/curve/proto/scan.proto

## heartbeat.proto
protoc --go_out=internal --proto_path=external/curve \
    --go_opt=Mproto/common.proto=github.com/opencurve/curve-manager/internal/proto/common \
    --go_opt=Mproto/common.proto=github.com/opencurve/curve-manager/internal/proto/scan \
    external/curve/proto/heartbeat.proto

## integrity.proto
protoc --go_out=internal --proto_path=external/curve \
    external/curve/proto/integrity.proto

## namespace2.proto
protoc --go_out=internal --proto_path=external/curve \
    --go_opt=Mproto/common.proto=github.com/opencurve/curve-manager/internal/proto/common \
    external/curve/proto/nameserver2.proto

## schedule.proto
protoc --go_out=internal --proto_path=external/curve \
    external/curve/proto/schedule.proto

## snapshotcloneserver.proto
protoc --go_out=internal --proto_path=external/curve \
    external/curve/proto/snapshotcloneserver.proto

## topology.proto
protoc --go_out=internal --proto_path=external/curve \
    --go_opt=Mproto/common.proto=github.com/opencurve/curve-manager/internal/proto/common \
    external/curve/proto/topology.proto

## statuscode.proto
protoc --go_out=internal --proto_path=external/curve/tools-v2 \
    external/curve/tools-v2/internal/proto/curvebs/topology/statuscode.proto

protoc --go-grpc_out=internal --proto_path=external/curve \
    external/curve/proto/*.proto