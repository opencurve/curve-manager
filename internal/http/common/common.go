package common

import (
	"fmt"
	"google.golang.org/protobuf/runtime/protoimpl"
	"strings"
)

const (
	GiB         = 1024 * 1024 * 1024
	TIME_FORMAT = "2006-01-02 15:04:05"
)

type CopysetInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LogicalPoolId      *uint32 `protobuf:"varint,1,req,name=logicalPoolId" json:"logicalPoolId,omitempty"`
	CopysetId          *uint32 `protobuf:"varint,2,req,name=copysetId" json:"copysetId,omitempty"`
	Scaning            *bool   `protobuf:"varint,3,opt,name=scaning" json:"scaning,omitempty"`
	LastScanSec        *uint64 `protobuf:"varint,4,opt,name=lastScanSec" json:"lastScanSec,omitempty"`
	LastScanConsistent *bool   `protobuf:"varint,5,opt,name=lastScanConsistent" json:"lastScanConsistent,omitempty"`
}

type ChunkServerLocation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChunkServerID *uint32 `protobuf:"varint,1,req,name=chunkServerID" json:"chunkServerID,omitempty"`
	HostIp        *string `protobuf:"bytes,2,req,name=hostIp" json:"hostIp,omitempty"`
	Port          *uint32 `protobuf:"varint,3,req,name=port" json:"port,omitempty"`
	ExternalIp    *string `protobuf:"bytes,4,opt,name=externalIp" json:"externalIp,omitempty"`
}

func (x *CopysetInfo) GetLogicalPoolId() uint32 {
	if x != nil && x.LogicalPoolId != nil {
		return *x.LogicalPoolId
	}
	return 0
}

func (x *CopysetInfo) GetCopysetId() uint32 {
	if x != nil && x.CopysetId != nil {
		return *x.CopysetId
	}
	return 0
}

func (x *CopysetInfo) GetScaning() bool {
	if x != nil && x.Scaning != nil {
		return *x.Scaning
	}
	return false
}

func (x *CopysetInfo) GetLastScanSec() uint64 {
	if x != nil && x.LastScanSec != nil {
		return *x.LastScanSec
	}
	return 0
}

func (x *CopysetInfo) GetLastScanConsistent() bool {
	if x != nil && x.LastScanConsistent != nil {
		return *x.LastScanConsistent
	}
	return false
}

func (x *ChunkServerLocation) GetChunkServerID() uint32 {
	if x != nil && x.ChunkServerID != nil {
		return *x.ChunkServerID
	}
	return 0
}

func (x *ChunkServerLocation) GetHostIp() string {
	if x != nil && x.HostIp != nil {
		return *x.HostIp
	}
	return ""
}

func (x *ChunkServerLocation) GetPort() uint32 {
	if x != nil && x.Port != nil {
		return *x.Port
	}
	return 0
}

func (x *ChunkServerLocation) GetExternalIp() string {
	if x != nil && x.ExternalIp != nil {
		return *x.ExternalIp
	}
	return ""
}

func ParseBvarMetric(value string) (*map[string]string, error) {
	ret := make(map[string]string)
	lines := strings.Split(value, "\n")
	for _, line := range lines {
		items := strings.Split(line, " : ")
		if len(items) != 2 {
			return nil, fmt.Errorf("parseBvarMetric failed, line: %s", line)
		}
		ret[strings.TrimSpace(items[0])] = strings.Trim(strings.TrimSpace(items[1]), "\"")
	}
	return &ret, nil
}
