package bsmetric

import (
	"testing"
)

func TestMdsStatus(t *testing.T) {
	c, err := GetMdsStatus()
	print(c)
	print(err)
}
