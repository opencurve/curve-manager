package nameserver2

import "google.golang.org/protobuf/runtime/protoimpl"

type FileType int32

const (
	FileType_INODE_DIRECTORY         FileType = 0
	FileType_INODE_PAGEFILE          FileType = 1
	FileType_INODE_APPENDFILE        FileType = 2
	FileType_INODE_APPENDECFILE      FileType = 3
	FileType_INODE_SNAPSHOT_PAGEFILE FileType = 4
)

var (
	StatusCode_name = map[int32]string{
		0:   "kOK",
		101: "kFileExists",
		102: "kFileNotExists",
		103: "kNotDirectory",
		104: "kParaError",
		105: "kShrinkBiggerFile",
		106: "kExtentUnitError",
		107: "kSegmentNotAllocated",
		108: "kSegmentAllocateError",
		109: "kDirNotExist",
		110: "kNotSupported",
		111: "kOwnerAuthFail",
		112: "kDirNotEmpty",
		120: "kFileUnderSnapShot",
		121: "kFileNotUnderSnapShot",
		122: "kSnapshotDeleting",
		123: "kSnapshotFileNotExists",
		124: "kSnapshotFileDeleteError",
		125: "kSessionNotExist",
		126: "kFileOccupied",
		127: "kCloneFileNameIllegal",
		128: "kCloneStatusNotMatch",
		129: "kCommonFileDeleteError",
		130: "kFileIdNotMatch",
		131: "kFileUnderDeleting",
		132: "kFileLengthNotSupported",
		133: "kDeleteFileBeingCloned",
		134: "kClientVersionNotMatch",
		135: "kSnapshotFrozen",
		136: "kSnapshotCloneConnectFail",
		137: "kSnapshotCloneServerNotInit",
		138: "kRecoverFileCloneMetaInstalled",
		139: "kRecoverFileError",
		140: "kEpochTooOld",
		501: "kStorageError",
		502: "KInternalError",
	}
	StatusCode_value = map[string]int32{
		"kOK":                            0,
		"kFileExists":                    101,
		"kFileNotExists":                 102,
		"kNotDirectory":                  103,
		"kParaError":                     104,
		"kShrinkBiggerFile":              105,
		"kExtentUnitError":               106,
		"kSegmentNotAllocated":           107,
		"kSegmentAllocateError":          108,
		"kDirNotExist":                   109,
		"kNotSupported":                  110,
		"kOwnerAuthFail":                 111,
		"kDirNotEmpty":                   112,
		"kFileUnderSnapShot":             120,
		"kFileNotUnderSnapShot":          121,
		"kSnapshotDeleting":              122,
		"kSnapshotFileNotExists":         123,
		"kSnapshotFileDeleteError":       124,
		"kSessionNotExist":               125,
		"kFileOccupied":                  126,
		"kCloneFileNameIllegal":          127,
		"kCloneStatusNotMatch":           128,
		"kCommonFileDeleteError":         129,
		"kFileIdNotMatch":                130,
		"kFileUnderDeleting":             131,
		"kFileLengthNotSupported":        132,
		"kDeleteFileBeingCloned":         133,
		"kClientVersionNotMatch":         134,
		"kSnapshotFrozen":                135,
		"kSnapshotCloneConnectFail":      136,
		"kSnapshotCloneServerNotInit":    137,
		"kRecoverFileCloneMetaInstalled": 138,
		"kRecoverFileError":              139,
		"kEpochTooOld":                   140,
		"kStorageError":                  501,
		"KInternalError":                 502,
	}
)

type StatusCode int32

const (
	// 执行成功
	StatusCode_kOK StatusCode = 0
	// 文件已存在
	StatusCode_kFileExists StatusCode = 101
	// 文件不存在
	StatusCode_kFileNotExists StatusCode = 102
	// 非目录类型
	StatusCode_kNotDirectory StatusCode = 103
	// 传入参数错误
	StatusCode_kParaError StatusCode = 104
	// 缩小文件，目前不支持缩小文件
	StatusCode_kShrinkBiggerFile StatusCode = 105
	// 扩容单位错误，非segment size整数倍
	StatusCode_kExtentUnitError StatusCode = 106
	// segment未分配
	StatusCode_kSegmentNotAllocated StatusCode = 107
	// segment分配失败
	StatusCode_kSegmentAllocateError StatusCode = 108
	// 目录不存在
	StatusCode_kDirNotExist StatusCode = 109
	// 功能不支持
	StatusCode_kNotSupported StatusCode = 110
	// owner认证失败
	StatusCode_kOwnerAuthFail StatusCode = 111
	// 目录非空
	StatusCode_kDirNotEmpty StatusCode = 112
	// 文件已处于快照中
	StatusCode_kFileUnderSnapShot StatusCode = 120
	// 文件不在快照中
	StatusCode_kFileNotUnderSnapShot StatusCode = 121
	// 快照删除中
	StatusCode_kSnapshotDeleting StatusCode = 122
	// 快照文件不存在
	StatusCode_kSnapshotFileNotExists StatusCode = 123
	// 快照文件删除失败
	StatusCode_kSnapshotFileDeleteError StatusCode = 124
	// session不存在
	StatusCode_kSessionNotExist StatusCode = 125
	// 文件已被占用
	StatusCode_kFileOccupied         StatusCode = 126
	StatusCode_kCloneFileNameIllegal StatusCode = 127
	StatusCode_kCloneStatusNotMatch  StatusCode = 128
	// 文件删除失败
	StatusCode_kCommonFileDeleteError StatusCode = 129
	// 文件id不匹配
	StatusCode_kFileIdNotMatch StatusCode = 130
	// 文件在删除中
	StatusCode_kFileUnderDeleting StatusCode = 131
	// 文件长度不符合要求
	StatusCode_kFileLengthNotSupported StatusCode = 132
	// 文件正在被克隆
	StatusCode_kDeleteFileBeingCloned StatusCode = 133
	// client版本不匹配
	StatusCode_kClientVersionNotMatch StatusCode = 134
	// snapshot功能禁用中
	StatusCode_kSnapshotFrozen StatusCode = 135
	// 快照克隆服务连不上
	StatusCode_kSnapshotCloneConnectFail StatusCode = 136
	// 快照克隆服务未初始化
	StatusCode_kSnapshotCloneServerNotInit StatusCode = 137
	// recover file status is CloneMetaInstalled
	StatusCode_kRecoverFileCloneMetaInstalled StatusCode = 138
	// recover file fail
	StatusCode_kRecoverFileError StatusCode = 139
	// epoch too old
	StatusCode_kEpochTooOld StatusCode = 140
	// 元数据存储错误
	StatusCode_kStorageError StatusCode = 501
	// 内部错误
	StatusCode_KInternalError StatusCode = 502
)

type GetAllocatedSizeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode *StatusCode `protobuf:"varint,1,req,name=statusCode,enum=curve.mds.StatusCode" json:"statusCode,omitempty"`
	// 文件或目录的分配大小
	AllocatedSize *uint64 `protobuf:"varint,2,opt,name=allocatedSize" json:"allocatedSize,omitempty"`
	// key是逻辑池id，value是分配大小
	AllocSizeMap map[uint32]uint64 `protobuf:"bytes,3,rep,name=allocSizeMap" json:"allocSizeMap,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
}

type FileStatus int32

const (
	// 文件创建完成
	FileStatus_kFileCreated FileStatus = 0
	// 文件删除中
	FileStatus_kFileDeleting FileStatus = 1
	// 文件正在克隆
	FileStatus_kFileCloning FileStatus = 2
	// 文件元数据安装完毕
	FileStatus_kFileCloneMetaInstalled FileStatus = 3
	// 文件克隆完成
	FileStatus_kFileCloned FileStatus = 4
	// 文件正在被克隆
	FileStatus_kFileBeingCloned FileStatus = 5
)

// Enum value maps for FileStatus.
var (
	FileStatus_name = map[int32]string{
		0: "kFileCreated",
		1: "kFileDeleting",
		2: "kFileCloning",
		3: "kFileCloneMetaInstalled",
		4: "kFileCloned",
		5: "kFileBeingCloned",
	}
	FileStatus_value = map[string]int32{
		"kFileCreated":            0,
		"kFileDeleting":           1,
		"kFileCloning":            2,
		"kFileCloneMetaInstalled": 3,
		"kFileCloned":             4,
		"kFileBeingCloned":        5,
	}
)

type ThrottleType int32

const (
	ThrottleType_IOPS_TOTAL ThrottleType = 1
	ThrottleType_IOPS_READ  ThrottleType = 2
	ThrottleType_IOPS_WRITE ThrottleType = 3
	ThrottleType_BPS_TOTAL  ThrottleType = 4
	ThrottleType_BPS_READ   ThrottleType = 5
	ThrottleType_BPS_WRITE  ThrottleType = 6
)

// Enum value maps for ThrottleType.
var (
	ThrottleType_name = map[int32]string{
		1: "IOPS_TOTAL",
		2: "IOPS_READ",
		3: "IOPS_WRITE",
		4: "BPS_TOTAL",
		5: "BPS_READ",
		6: "BPS_WRITE",
	}
	ThrottleType_value = map[string]int32{
		"IOPS_TOTAL": 1,
		"IOPS_READ":  2,
		"IOPS_WRITE": 3,
		"BPS_TOTAL":  4,
		"BPS_READ":   5,
		"BPS_WRITE":  6,
	}
)

type ListDirResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode *StatusCode `protobuf:"varint,1,req,name=statusCode,enum=curve.mds.StatusCode" json:"statusCode,omitempty"`
	FileInfo   []*FileInfo `protobuf:"bytes,2,rep,name=fileInfo" json:"fileInfo,omitempty"`
}

func (x *GetAllocatedSizeResponse) GetStatusCode() StatusCode {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return StatusCode_kOK
}

func (x *GetAllocatedSizeResponse) GetAllocatedSize() uint64 {
	if x != nil && x.AllocatedSize != nil {
		return *x.AllocatedSize
	}
	return 0
}

func (x *GetAllocatedSizeResponse) GetAllocSizeMap() map[uint32]uint64 {
	if x != nil {
		return x.AllocSizeMap
	}
	return nil
}

func (x *ListDirResponse) GetStatusCode() StatusCode {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return StatusCode_kOK
}

func (x *ListDirResponse) GetFileInfo() []*FileInfo {
	if x != nil {
		return x.FileInfo
	}
	return nil
}

type FileInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          *uint64     `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	FileName    *string     `protobuf:"bytes,2,opt,name=fileName" json:"fileName,omitempty"`
	ParentId    *uint64     `protobuf:"varint,3,opt,name=parentId" json:"parentId,omitempty"`
	FileType    *FileType   `protobuf:"varint,4,opt,name=fileType,enum=curve.mds.FileType" json:"fileType,omitempty"`
	Owner       *string     `protobuf:"bytes,5,opt,name=owner" json:"owner,omitempty"`
	ChunkSize   *uint32     `protobuf:"varint,6,opt,name=chunkSize" json:"chunkSize,omitempty"`
	SegmentSize *uint32     `protobuf:"varint,7,opt,name=segmentSize" json:"segmentSize,omitempty"`
	Length      *uint64     `protobuf:"varint,8,opt,name=length" json:"length,omitempty"`
	Ctime       *uint64     `protobuf:"varint,9,opt,name=ctime" json:"ctime,omitempty"`
	SeqNum      *uint64     `protobuf:"varint,10,opt,name=seqNum" json:"seqNum,omitempty"`
	FileStatus  *FileStatus `protobuf:"varint,11,opt,name=fileStatus,enum=curve.mds.FileStatus" json:"fileStatus,omitempty"`
	// 用于文件转移到回收站的情况下恢复场景下的使用,
	// RecycleBin（回收站）目录下使用/其他场景下不使用
	OriginalFullPathName *string `protobuf:"bytes,12,opt,name=originalFullPathName" json:"originalFullPathName,omitempty"`
	// cloneSource 当前用于存放克隆源(当前主要用于curvefs)
	// 后期可以考虑存放 s3相关信息
	CloneSource *string `protobuf:"bytes,13,opt,name=cloneSource" json:"cloneSource,omitempty"`
	// cloneLength 克隆源文件的长度，用于clone过程中进行extent
	CloneLength    *uint64             `protobuf:"varint,14,opt,name=cloneLength" json:"cloneLength,omitempty"`
	StripeUnit     *uint64             `protobuf:"varint,15,opt,name=stripeUnit" json:"stripeUnit,omitempty"`
	StripeCount    *uint64             `protobuf:"varint,16,opt,name=stripeCount" json:"stripeCount,omitempty"`
	ThrottleParams *FileThrottleParams `protobuf:"bytes,17,opt,name=throttleParams" json:"throttleParams,omitempty"`
	Epoch          *uint64             `protobuf:"varint,18,opt,name=epoch" json:"epoch,omitempty"`
}

func (x *FileInfo) GetId() uint64 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return 0
}

func (x *FileInfo) GetFileName() string {
	if x != nil && x.FileName != nil {
		return *x.FileName
	}
	return ""
}

func (x *FileInfo) GetParentId() uint64 {
	if x != nil && x.ParentId != nil {
		return *x.ParentId
	}
	return 0
}

func (x *FileInfo) GetFileType() FileType {
	if x != nil && x.FileType != nil {
		return *x.FileType
	}
	return FileType_INODE_DIRECTORY
}

func (x *FileInfo) GetOwner() string {
	if x != nil && x.Owner != nil {
		return *x.Owner
	}
	return ""
}

func (x *FileInfo) GetChunkSize() uint32 {
	if x != nil && x.ChunkSize != nil {
		return *x.ChunkSize
	}
	return 0
}

func (x *FileInfo) GetSegmentSize() uint32 {
	if x != nil && x.SegmentSize != nil {
		return *x.SegmentSize
	}
	return 0
}

func (x *FileInfo) GetLength() uint64 {
	if x != nil && x.Length != nil {
		return *x.Length
	}
	return 0
}

func (x *FileInfo) GetCtime() uint64 {
	if x != nil && x.Ctime != nil {
		return *x.Ctime
	}
	return 0
}

func (x *FileInfo) GetSeqNum() uint64 {
	if x != nil && x.SeqNum != nil {
		return *x.SeqNum
	}
	return 0
}

func (x *FileInfo) GetFileStatus() FileStatus {
	if x != nil && x.FileStatus != nil {
		return *x.FileStatus
	}
	return FileStatus_kFileCreated
}

func (x *FileInfo) GetOriginalFullPathName() string {
	if x != nil && x.OriginalFullPathName != nil {
		return *x.OriginalFullPathName
	}
	return ""
}

func (x *FileInfo) GetCloneSource() string {
	if x != nil && x.CloneSource != nil {
		return *x.CloneSource
	}
	return ""
}

func (x *FileInfo) GetCloneLength() uint64 {
	if x != nil && x.CloneLength != nil {
		return *x.CloneLength
	}
	return 0
}

func (x *FileInfo) GetStripeUnit() uint64 {
	if x != nil && x.StripeUnit != nil {
		return *x.StripeUnit
	}
	return 0
}

func (x *FileInfo) GetStripeCount() uint64 {
	if x != nil && x.StripeCount != nil {
		return *x.StripeCount
	}
	return 0
}

func (x *FileInfo) GetThrottleParams() *FileThrottleParams {
	if x != nil {
		return x.ThrottleParams
	}
	return nil
}

func (x *FileInfo) GetEpoch() uint64 {
	if x != nil && x.Epoch != nil {
		return *x.Epoch
	}
	return 0
}

type FileThrottleParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ThrottleParams []*ThrottleParams `protobuf:"bytes,1,rep,name=throttleParams" json:"throttleParams,omitempty"`
}

func (x *FileThrottleParams) GetThrottleParams() []*ThrottleParams {
	if x != nil {
		return x.ThrottleParams
	}
	return nil
}

type ThrottleParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type        *ThrottleType `protobuf:"varint,1,req,name=type,enum=curve.mds.ThrottleType" json:"type,omitempty"`
	Limit       *uint64       `protobuf:"varint,2,req,name=limit" json:"limit,omitempty"`
	Burst       *uint64       `protobuf:"varint,3,opt,name=burst" json:"burst,omitempty"`
	BurstLength *uint64       `protobuf:"varint,4,opt,name=burstLength" json:"burstLength,omitempty"`
}

func (x *ThrottleParams) GetType() ThrottleType {
	if x != nil && x.Type != nil {
		return *x.Type
	}
	return ThrottleType_IOPS_TOTAL
}

func (x *ThrottleParams) GetLimit() uint64 {
	if x != nil && x.Limit != nil {
		return *x.Limit
	}
	return 0
}

func (x *ThrottleParams) GetBurst() uint64 {
	if x != nil && x.Burst != nil {
		return *x.Burst
	}
	return 0
}

func (x *ThrottleParams) GetBurstLength() uint64 {
	if x != nil && x.BurstLength != nil {
		return *x.BurstLength
	}
	return 0
}

type GetFileInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode *StatusCode `protobuf:"varint,1,req,name=statusCode,enum=curve.mds.StatusCode" json:"statusCode,omitempty"`
	FileInfo   *FileInfo   `protobuf:"bytes,2,opt,name=fileInfo" json:"fileInfo,omitempty"`
}

type GetFileSizeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode *StatusCode `protobuf:"varint,1,req,name=statusCode,enum=curve.mds.StatusCode" json:"statusCode,omitempty"`
	// 文件或目录的file length
	FileSize *uint64 `protobuf:"varint,2,opt,name=fileSize" json:"fileSize,omitempty"`
}

func (x *GetFileSizeResponse) GetStatusCode() StatusCode {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return StatusCode_kOK
}

func (x *GetFileSizeResponse) GetFileSize() uint64 {
	if x != nil && x.FileSize != nil {
		return *x.FileSize
	}
	return 0
}

func (x *GetFileInfoResponse) GetStatusCode() StatusCode {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return StatusCode_kOK
}

func (x *GetFileInfoResponse) GetFileInfo() *FileInfo {
	if x != nil {
		return x.FileInfo
	}
	return nil
}

type DeleteFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode *StatusCode `protobuf:"varint,1,req,name=statusCode,enum=curve.mds.StatusCode" json:"statusCode,omitempty"`
}
type RecoverFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode *StatusCode `protobuf:"varint,1,req,name=statusCode,enum=curve.mds.StatusCode" json:"statusCode,omitempty"`
}

func (x *DeleteFileResponse) GetStatusCode() StatusCode {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return StatusCode_kOK
}

func (x *RecoverFileResponse) GetStatusCode() StatusCode {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return StatusCode_kOK
}

type CreateFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode *StatusCode `protobuf:"varint,1,req,name=statusCode,enum=curve.mds.StatusCode" json:"statusCode,omitempty"`
}

func (x *CreateFileResponse) GetStatusCode() StatusCode {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return StatusCode_kOK
}

type ExtendFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode *StatusCode `protobuf:"varint,1,req,name=statusCode,enum=curve.mds.StatusCode" json:"statusCode,omitempty"`
}

func (x *ExtendFileResponse) GetStatusCode() StatusCode {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return StatusCode_kOK
}

type UpdateFileThrottleParamsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode *StatusCode `protobuf:"varint,1,req,name=statusCode,enum=curve.mds.StatusCode" json:"statusCode,omitempty"`
}

func (x *UpdateFileThrottleParamsResponse) GetStatusCode() StatusCode {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return StatusCode_kOK
}

type FindFileMountPointResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode *StatusCode   `protobuf:"varint,1,req,name=statusCode,enum=curve.mds.StatusCode" json:"statusCode,omitempty"`
	ClientInfo []*ClientInfo `protobuf:"bytes,2,rep,name=clientInfo" json:"clientInfo,omitempty"`
}

func (x *FindFileMountPointResponse) GetStatusCode() StatusCode {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return StatusCode_kOK
}

func (x *FindFileMountPointResponse) GetClientInfo() []*ClientInfo {
	if x != nil {
		return x.ClientInfo
	}
	return nil
}

type ClientInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip   *string `protobuf:"bytes,1,req,name=ip" json:"ip,omitempty"`
	Port *uint32 `protobuf:"varint,2,req,name=port" json:"port,omitempty"`
}

func (x *ClientInfo) GetIp() string {
	if x != nil && x.Ip != nil {
		return *x.Ip
	}
	return ""
}

func (x *ClientInfo) GetPort() uint32 {
	if x != nil && x.Port != nil {
		return *x.Port
	}
	return 0
}
