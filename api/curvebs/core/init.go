package core

import (
	"github.com/opencurve/curve-manager/internal/email"
	metrics "github.com/opencurve/curve-manager/internal/metrics/core"
	bsrpc "github.com/opencurve/curve-manager/internal/rpc/curvebs"
	"github.com/opencurve/curve-manager/internal/snapshotclone"
	"github.com/opencurve/curve-manager/internal/storage"
	"github.com/opencurve/pigeon"
)

func initClients(cfg *pigeon.Configure) error {
	// init mds rpc client
	err := bsrpc.Init(cfg)
	if err != nil {
		return err
	}

	// init metric client
	err = metrics.Init(cfg)
	if err != nil {
		return err
	}

	// init snapshot clone client
	snapshotclone.Init(cfg)

	// init email client which used to reset password and some system notifications
	email.Init(cfg)
	return nil
}

func Init(cfg *pigeon.Configure) error {
	// init storage
	err := storage.Init(cfg)
	if err != nil {
		return err
	}

	return initClients(cfg)
}
