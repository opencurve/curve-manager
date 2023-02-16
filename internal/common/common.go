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

package common

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	GiB                       = 1024 * 1024 * 1024
	TIME_FORMAT               = "2006-01-02 15:04:05"
	CURVEBS_ADDRESS_DELIMITER = ","
	RAFT_REPLICAS_NUMBER      = 3
	RAFT_MARGIN               = 1000

	CHAR_TABLE = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// raft status
	RAFT_EMPTY_ADDR                            = "0.0.0.0:0:0"
	RAFT_STATUS_KEY_GROUPID                    = "group_id"
	RAFT_STATUS_KEY_LEADER                     = "leader"
	RAFT_STATUS_KEY_PEERS                      = "peers"
	RAFT_STATUS_KEY_STATE                      = "state"
	RAFT_STATUS_KEY_REPLICATOR                 = "replicator"
	RAFT_STATUS_KEY_LAST_LOG_ID                = "last_log_id"
	RAFT_STATUS_KEY_SNAPSHOT                   = "snapshot"
	RAFT_STATUS_KEY_NEXT_INDEX                 = "next_index"
	RAFT_STATUS_KEY_FLYING_APPEND_ENTRIES_SIZE = "flying_append_entries_size"
	RAFT_STATUS_KEY_STORAGE                    = "storage"

	RAFT_STATUS_STATE_LEADER       = "LEADER"
	RAFT_STATUS_STATE_FOLLOWER     = "FOLLOWER"
	RAFT_STATUS_STATE_TRANSFERRING = "TRANSFERRING"
	RAFT_STATUS_STATE_CANDIDATE    = "CANDIDATE"
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

func Max(first, second uint64) uint64 {
	if first < second {
		return second
	}
	return first
}

func GetMd5Sum32Little(key string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}

func GetRandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = CHAR_TABLE[rand.Intn(len(CHAR_TABLE))]
	}
	return string(b)
}

func GetIPFromEndpoint(endpoint string) (string, error) {
	strs := strings.Split(endpoint, ":")
	if len(strs) != 2 {
		return "", fmt.Errorf("invalid endpoint")
	}
	return strs[0], nil
}