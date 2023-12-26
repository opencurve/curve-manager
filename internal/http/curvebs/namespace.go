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
* Project: Curve-Go-RPC
* Created Date: 2023-03-03
* Author: wanghai (SeanHai)
 */

package curvebs

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/opencurve/curve-manager/internal/http/common"
	"github.com/opencurve/curve-manager/internal/http/nameserver2"
	"strconv"
	"time"
)

const (
	// file type
	INODE_DIRECTORY         = "INODE_DIRECTORY"
	INODE_PAGEFILE          = "INODE_PAGEFILE"
	INODE_APPENDFILE        = "INODE_APPENDFILE"
	INODE_APPENDECFILE      = "INODE_APPENDECFILE"
	INODE_SNAPSHOT_PAGEFILE = "INODE_SNAPSHOT_PAGEFILE"

	// file status
	FILE_CREATED             = "kFileCreated"
	FILE_DELETING            = "kFileDeleting"
	FILE_CLONING             = "kFileCloning"
	FILE_CLONEMETA_INSTALLED = "kFileCloneMetaInstalled"
	FILE_CLONED              = "kFileCloned"
	FILE_BEIING_CLONED       = "kFileBeingCloned"

	// throttle type
	IOPS_TOTAL = "IOPS_TOTAL"
	IOPS_READ  = "IOPS_READ"
	IOPS_WRITE = "IOPS_WRITE"
	BPS_TOTAL  = "BPS_TOTAL"
	BPS_READ   = "BPS_READ"
	BPS_WRITE  = "BPS_WRITE"

	// apis
	GET_FILE_ALLOC_SIZE_FUNC    = "GetAllocatedSize"
	LIST_DIR_FUNC               = "ListDir"
	GET_FILE_INFO               = "GetFileInfo"
	GET_FILE_SIZE               = "GetFileSize"
	DELETE_FILE                 = "DeleteFile"
	CREATE_FILE                 = "CreateFile"
	EXTEND_FILE                 = "ExtendFile"
	RECOVER_FILE                = "RecoverFile"
	UPDATE_FILE_THROTTLE_PARAMS = "UpdateFileThrottleParams"
	FIND_FILE_MOUNTPOINT        = "FindFileMountPoint"

	GET_FILE_ALLOC_SIZE_FUNC_http    = "GetAllocatedSize"
	LIST_DIR_FUNC_http               = "ListDir"
	GET_FILE_INFO_http               = "GetFileInfo"
	GET_FILE_SIZE_http               = "GetFileSize"
	DELETE_FILE_http                 = "DeleteFile"
	CREATE_FILE_http                 = "CreateFile"
	EXTEND_FILE_http                 = "ExtendFile"
	RECOVER_FILE_http                = "RecoverFile"
	UPDATE_FILE_THROTTLE_PARAMS_http = "UpdateFileThrottleParams"
	FIND_FILE_MOUNTPOINT_http        = "FindFileMountPoint"
)

type ThrottleParams struct {
	Type        string `json:"type"`
	Limit       uint64 `json:"limit"`
	Burst       uint64 `json:"burst"`
	BurstLength uint64 `json:"burstLength"`
}

type FileInfo struct {
	Id                   uint64           `json:"id"`
	FileName             string           `json:"fileName"`
	ParentId             uint64           `json:"parentId"`
	FileType             string           `json:"fileType"`
	Owner                string           `json:"owner"`
	ChunkSize            uint32           `json:"chunkSize"`
	SegmentSize          uint32           `json:"segmentSize"`
	Length               uint64           `json:"length"`
	AllocateSize         uint64           `json:"alloc"`
	Ctime                string           `json:"ctime"`
	SeqNum               uint64           `json:"seqNum"`
	FileStatus           string           `json:"fileStatus"`
	OriginalFullPathName string           `json:"originalFullPathName"`
	CloneSource          string           `json:"cloneSource"`
	CloneLength          uint64           `json:"cloneLength"`
	StripeUnit           uint64           `json:"stripeUnit"`
	StripeCount          uint64           `json:"stripeCount"`
	ThrottleParams       []ThrottleParams `json:"throttleParams"`
	Epoch                uint64           `json:"epoch"`
	MountPoints          []string         `json:"mountPoints"`
}

func (cli *MdsClient) GetFileAllocatedSize(filename string) (uint64, map[uint32]uint64, error) {
	var host = cli.addrs
	//todo checkHost
	var path = GET_FILE_ALLOC_SIZE_FUNC_http
	path = path + "FileName=" + filename

	ret := cli.baseClient_http.SendHTTP(host, path)
	if ret.Err != nil {
		return 0, nil, ret.Err
	}
	v := ret.Result.(*resty.Response).String()
	var response *nameserver2.GetAllocatedSizeResponse
	err := json.Unmarshal([]byte(v), &response)
	if err != nil {
		return 0, nil, err
	}
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return 0, nil, fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	infos := make(map[uint32]uint64)
	for k, v := range response.GetAllocSizeMap() {
		infos[k] = v / common.GiB
	}
	return response.GetAllocatedSize() / common.GiB, infos, nil
}

func getFileType(t string) nameserver2.FileType {
	switch t {
	case INODE_DIRECTORY:
		return nameserver2.FileType_INODE_DIRECTORY
	case INODE_PAGEFILE:
		return nameserver2.FileType_INODE_PAGEFILE
	case INODE_APPENDFILE:
		return nameserver2.FileType_INODE_APPENDFILE
	case INODE_APPENDECFILE:
		return nameserver2.FileType_INODE_APPENDECFILE
	case INODE_SNAPSHOT_PAGEFILE:
		return nameserver2.FileType_INODE_SNAPSHOT_PAGEFILE
	default:
		return -1
	}
}

func getFileTypeStr(t nameserver2.FileType) string {
	switch t {
	case nameserver2.FileType_INODE_DIRECTORY:
		return INODE_DIRECTORY
	case nameserver2.FileType_INODE_PAGEFILE:
		return INODE_PAGEFILE
	case nameserver2.FileType_INODE_APPENDFILE:
		return INODE_APPENDFILE
	case nameserver2.FileType_INODE_APPENDECFILE:
		return INODE_APPENDECFILE
	case nameserver2.FileType_INODE_SNAPSHOT_PAGEFILE:
		return INODE_SNAPSHOT_PAGEFILE
	default:
		return INVALID
	}
}

func getFileStatus(s nameserver2.FileStatus) string {
	switch s {
	case nameserver2.FileStatus_kFileCreated:
		return FILE_CREATED
	case nameserver2.FileStatus_kFileDeleting:
		return FILE_DELETING
	case nameserver2.FileStatus_kFileCloning:
		return FILE_CLONING
	case nameserver2.FileStatus_kFileCloneMetaInstalled:
		return FILE_CLONEMETA_INSTALLED
	case nameserver2.FileStatus_kFileCloned:
		return FILE_CLONED
	case nameserver2.FileStatus_kFileBeingCloned:
		return FILE_BEIING_CLONED
	default:
		return INVALID
	}
}

func getThrottleTypeStr(t nameserver2.ThrottleType) string {
	switch t {
	case nameserver2.ThrottleType_IOPS_TOTAL:
		return IOPS_TOTAL
	case nameserver2.ThrottleType_IOPS_READ:
		return IOPS_READ
	case nameserver2.ThrottleType_IOPS_WRITE:
		return IOPS_WRITE
	case nameserver2.ThrottleType_BPS_TOTAL:
		return BPS_TOTAL
	case nameserver2.ThrottleType_BPS_READ:
		return BPS_READ
	case nameserver2.ThrottleType_BPS_WRITE:
		return BPS_WRITE
	default:
		return INVALID
	}
}

func getThrottleType(t string) nameserver2.ThrottleType {
	switch t {
	case IOPS_TOTAL:
		return nameserver2.ThrottleType_IOPS_TOTAL
	case IOPS_READ:
		return nameserver2.ThrottleType_IOPS_READ
	case IOPS_WRITE:
		return nameserver2.ThrottleType_IOPS_WRITE
	case BPS_TOTAL:
		return nameserver2.ThrottleType_BPS_TOTAL
	case BPS_READ:
		return nameserver2.ThrottleType_BPS_READ
	case BPS_WRITE:
		return nameserver2.ThrottleType_BPS_WRITE
	default:
		return 0
	}
}

func (cli *MdsClient) ListDir(filename, owner, sig string, date uint64) ([]FileInfo, error) {
	var host = cli.addrs
	//todo check URL
	var path = LIST_DIR_FUNC_http
	path = path + "FileName=" + filename + "&Owner=" + owner + "&Date=" + strconv.Itoa(int(date))
	if sig != "" {
		path = fmt.Sprintf("%s %s %d", path, "Signature=", &sig)
	}

	ret := cli.baseClient_http.SendHTTP(host, path)
	if ret.Err != nil {
		return nil, ret.Err
	}
	v := ret.Result.(*resty.Response).String()
	var response *nameserver2.ListDirResponse
	err := json.Unmarshal([]byte(v), &response)
	if err != nil {
		return nil, err
	}
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return nil, fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	infos := []FileInfo{}
	for _, v := range response.GetFileInfo() {
		var info FileInfo
		info.Id = v.GetId()
		info.FileName = v.GetFileName()
		info.ParentId = v.GetParentId()
		info.FileType = getFileTypeStr(v.GetFileType())
		info.Owner = v.GetOwner()
		info.ChunkSize = v.GetChunkSize()
		info.SegmentSize = v.GetSegmentSize()
		info.Length = v.GetLength() / common.GiB
		info.Ctime = time.Unix(int64(v.GetCtime()/1000000), 0).Format(common.TIME_FORMAT)
		info.SeqNum = v.GetSeqNum()
		info.FileStatus = getFileStatus(v.GetFileStatus())
		info.OriginalFullPathName = v.GetOriginalFullPathName()
		info.CloneSource = v.GetCloneSource()
		info.CloneLength = v.GetCloneLength()
		info.StripeUnit = v.GetStripeUnit()
		info.StripeCount = v.GetStripeCount()
		info.ThrottleParams = []ThrottleParams{}
		for _, p := range v.GetThrottleParams().GetThrottleParams() {
			var param ThrottleParams
			param.Type = getThrottleTypeStr(p.GetType())
			param.Limit = p.GetLimit()
			param.Burst = p.GetBurst()
			param.BurstLength = p.GetBurstLength()
			info.ThrottleParams = append(info.ThrottleParams, param)
		}
		info.Epoch = v.GetEpoch()
		infos = append(infos, info)
	}
	return infos, nil
}

func (cli *MdsClient) GetFileInfo(filename, owner, sig string, date uint64) (FileInfo, error) {
	info := FileInfo{}
	var host = cli.addrs
	// todo check URL
	var path = GET_FILE_SIZE_http
	path = path + "FileName=" + filename + "&Owner=" + owner + "&Date=" + strconv.Itoa(int(date))
	if sig != "" {
		path = fmt.Sprintf("%s %s %d", path, "Signature=", &sig)
	}

	ret := cli.baseClient_http.SendHTTP(host, path)
	if ret.Err != nil {
		return info, ret.Err
	}
	tmp := ret.Result.(*resty.Response).String()
	var response *nameserver2.GetFileInfoResponse
	err := json.Unmarshal([]byte(tmp), &response)
	if err != nil {
		return info, err
	}
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return info, fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	v := response.GetFileInfo()
	info.Id = v.GetId()
	info.FileName = v.GetFileName()
	info.ParentId = v.GetParentId()
	info.FileType = getFileTypeStr(v.GetFileType())
	info.Owner = v.GetOwner()
	info.ChunkSize = v.GetChunkSize()
	info.SegmentSize = v.GetSegmentSize()
	info.Length = v.GetLength() / common.GiB
	info.Ctime = time.Unix(int64(v.GetCtime()/1000000), 0).Format(common.TIME_FORMAT)
	info.SeqNum = v.GetSeqNum()
	info.FileStatus = getFileStatus(v.GetFileStatus())
	info.OriginalFullPathName = v.GetOriginalFullPathName()
	info.CloneSource = v.GetCloneSource()
	info.CloneLength = v.GetCloneLength()
	info.StripeUnit = v.GetStripeUnit()
	info.StripeCount = v.GetStripeCount()
	info.ThrottleParams = []ThrottleParams{}
	for _, p := range v.GetThrottleParams().GetThrottleParams() {
		var param ThrottleParams
		param.Type = getThrottleTypeStr(p.GetType())
		param.Limit = p.GetLimit()
		param.Burst = p.GetBurst()
		param.BurstLength = p.GetBurstLength()
		info.ThrottleParams = append(info.ThrottleParams, param)
	}
	info.Epoch = v.GetEpoch()
	return info, nil
}

func (cli *MdsClient) GetFileSize(fileName string) (uint64, error) {
	var size uint64
	var host = cli.addrs
	var path = GET_FILE_SIZE_http
	// todo checkURL
	path = fmt.Sprintf("%s %s %s", path, "FileName=", fileName)
	ret := cli.baseClient_http.SendHTTP(host, path)
	if ret.Err != nil {
		return size, ret.Err
	}

	tmp := ret.Result.(*resty.Response).String()
	var response *nameserver2.GetFileSizeResponse
	err := json.Unmarshal([]byte(tmp), &response)
	if err != nil {
		return 0, err
	}
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return size, fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	size = response.GetFileSize() / common.GiB
	return size, nil
}

func (cli *MdsClient) DeleteFile(filename, owner, sig string, fileId, date uint64, forceDelete bool) error {
	var host = cli.addrs
	var path = DELETE_FILE_http
	//todo checkURL
	path = path + "FileName=" + filename + "&Date=" + strconv.Itoa(int(date)) + "&ForceDelete=" + strconv.FormatBool(forceDelete)
	if sig != "" {
		path = fmt.Sprintf("%s %s %s", path, "Signature=", sig)
	}
	if fileId != 0 {
		path = fmt.Sprintf("%s %s %d", path, "FileId=", &fileId)
	}

	ret := cli.baseClient_http.SendHTTP(host, path)

	if ret.Err != nil {
		return ret.Err
	}
	tmp := ret.Result.(*resty.Response).String()
	var response *nameserver2.DeleteFileResponse
	err := json.Unmarshal([]byte(tmp), &response)
	if err != nil {
		return err
	}
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	return nil
}

func (cli *MdsClient) RecoverFile(filename, owner, sig string, fileId, date uint64) error {
	var host = cli.addrs
	var path = RECOVER_FILE_http
	//todo checkURL
	path = path + "FileName=" + filename + "&Owner=" + owner + "&Date" + strconv.Itoa(int(date))
	if sig != "" {
		path = fmt.Sprintf("%s %s %s", path, "?Signature=", sig)
	}
	if fileId != 0 {
		path = fmt.Sprintf("%s %s %d", path, "?FileId=", &fileId)
	}

	ret := cli.baseClient_http.SendHTTP(host, path)
	if ret.Err != nil {
		return ret.Err
	}
	tmp := ret.Result.(*resty.Response).String()
	var response *nameserver2.RecoverFileResponse
	err := json.Unmarshal([]byte(tmp), &response)
	if err != nil {
		return err
	}
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	return nil
}

func (cli *MdsClient) CreateFile(filename, ftype, owner, sig string, length, date, stripeUnit, stripeCount uint64) error {
	var host = cli.addrs
	var path = CREATE_FILE_http
	//todo: generating param
	ret := cli.baseClient_http.SendHTTP(host, path)
	if ret.Err != nil {
		return ret.Err
	}
	tmp := ret.Result.(*resty.Response).String()
	var response *nameserver2.CreateFileResponse
	err := json.Unmarshal([]byte(tmp), &response)
	if err != nil {
		return err
	}
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	return nil
}

func (cli *MdsClient) ExtendFile(filename, owner, sig string, newSize, date uint64) error {
	var host = cli.addrs
	var path = EXTEND_FILE_http
	//todo: checkURL
	path = path + "FileName=" + filename + "&NewSize=" + strconv.Itoa(int(newSize)) + "&Owner=" + owner + "&Date" + strconv.Itoa(int(date))
	if sig != "" {
		path = fmt.Sprintf("%s %s %s", path, "Signature=", sig)
	}
	ret := cli.baseClient_http.SendHTTP(host, path)
	if ret.Err != nil {
		return ret.Err
	}
	tmp := ret.Result.(*resty.Response).String()
	var response *nameserver2.ExtendFileResponse
	err := json.Unmarshal([]byte(tmp), &response)
	if err != nil {
		return err
	}
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	return nil
}

func (cli *MdsClient) UpdateFileThrottleParams(filename, owner, sig string, date uint64, params ThrottleParams) error {
	var host = cli.addrs
	var path = UPDATE_FILE_THROTTLE_PARAMS_http

	//todo : generating paramt
	if sig != "" {
		path = fmt.Sprintf("%s %s %s", path, "Signature=", sig)
	}
	ret := cli.baseClient_http.SendHTTP(host, path)
	if ret.Err != nil {
		return ret.Err
	}
	tmp := ret.Result.(*resty.Response).String()
	var response *nameserver2.UpdateFileThrottleParamsResponse
	err := json.Unmarshal([]byte(tmp), &response)
	if err != nil {
		return err
	}
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	return nil
}

func (cli *MdsClient) FindFileMountPoint(filename string) ([]string, error) {
	info := []string{}

	var host = cli.addrs
	//todo: checkURL
	var path = FIND_FILE_MOUNTPOINT_http
	path = fmt.Sprintf("%s %s %s", path, "FileName=", filename)
	ret := cli.baseClient_http.SendHTTP(host, path)

	if ret.Err != nil {
		return nil, ret.Err
	}
	tmp := ret.Result.(*resty.Response).String()
	var response *nameserver2.FindFileMountPointResponse
	err := json.Unmarshal([]byte(tmp), &response)
	if err != nil {
		return nil, err
	}
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return info, fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	for _, v := range response.GetClientInfo() {
		info = append(info, fmt.Sprintf("%s:%d", v.GetIp(), v.GetPort()))
	}
	return info, nil
}
