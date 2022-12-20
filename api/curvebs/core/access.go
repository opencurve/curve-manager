package core

import (
	"github.com/opencurve/pigeon"
)

func AccessAllowed(r *pigeon.Request, data interface{}) bool {
	return true
}
