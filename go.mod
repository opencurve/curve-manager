module github.com/opencurve/curve-manager

go 1.16

require (
	github.com/go-resty/resty/v2 v2.7.0
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/mattn/go-sqlite3 v1.6.0
	github.com/opencurve/pigeon v0.6.0
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.1
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.49.0

replace github.com/opencurve/pigeon => ./external/pigeon
