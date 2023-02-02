package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const (
	GB                        = 1024 * 1024 * 1024
	TIME_FORMAT               = "2006-01-02 15:04:05"
	CURVEBS_ADDRESS_DELIMITER = ","
)

type QueryResult struct {
	Key    interface{}
	Err    error
	Result interface{}
}

func GetString2Signature(date int64, owner string) string {
	return fmt.Sprintf("%d:%s", date, owner)
}

func CalcString2Signature(in string, secretKet string) string {
	h := hmac.New(sha256.New, []byte(secretKet))
	h.Write([]byte(in))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
