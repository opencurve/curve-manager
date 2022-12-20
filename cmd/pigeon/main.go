package main

import (
	"github.com/opencurve/curve-manager/api/curvebs"
	"github.com/opencurve/pigeon"
)

func main() {
	curvebsServer := curvebs.NewServer()
	pigeon.Serve(curvebsServer)
}
