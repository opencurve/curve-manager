package statuscode

type TopoStatusCode int32

const (
	TopoStatusCode_Success                            TopoStatusCode = 0
	TopoStatusCode_InternalError                      TopoStatusCode = -1
	TopoStatusCode_InvalidParam                       TopoStatusCode = -2
	TopoStatusCode_InitFail                           TopoStatusCode = -3
	TopoStatusCode_StorgeFail                         TopoStatusCode = -4
	TopoStatusCode_IdDuplicated                       TopoStatusCode = -5
	TopoStatusCode_ChunkServerNotFound                TopoStatusCode = -6
	TopoStatusCode_ServerNotFound                     TopoStatusCode = -7
	TopoStatusCode_ZoneNotFound                       TopoStatusCode = -8
	TopoStatusCode_PhysicalPoolNotFound               TopoStatusCode = -9
	TopoStatusCode_LogicalPoolNotFound                TopoStatusCode = -10
	TopoStatusCode_CopySetNotFound                    TopoStatusCode = -11
	TopoStatusCode_GenCopysetErr                      TopoStatusCode = -12
	TopoStatusCode_AllocateIdFail                     TopoStatusCode = -13
	TopoStatusCode_CannotRemoveWhenNotEmpty           TopoStatusCode = -14
	TopoStatusCode_IpPortDuplicated                   TopoStatusCode = -15
	TopoStatusCode_NameDuplicated                     TopoStatusCode = -16
	TopoStatusCode_CreateCopysetNodeOnChunkServerFail TopoStatusCode = -17
	TopoStatusCode_CannotRemoveNotRetired             TopoStatusCode = -18
	TopoStatusCode_LogicalPoolExist                   TopoStatusCode = -19
)

// Enum value maps for TopoStatusCode.
var (
	TopoStatusCode_name = map[int32]string{
		0:   "Success",
		-1:  "InternalError",
		-2:  "InvalidParam",
		-3:  "InitFail",
		-4:  "StorgeFail",
		-5:  "IdDuplicated",
		-6:  "ChunkServerNotFound",
		-7:  "ServerNotFound",
		-8:  "ZoneNotFound",
		-9:  "PhysicalPoolNotFound",
		-10: "LogicalPoolNotFound",
		-11: "CopySetNotFound",
		-12: "GenCopysetErr",
		-13: "AllocateIdFail",
		-14: "CannotRemoveWhenNotEmpty",
		-15: "IpPortDuplicated",
		-16: "NameDuplicated",
		-17: "CreateCopysetNodeOnChunkServerFail",
		-18: "CannotRemoveNotRetired",
		-19: "LogicalPoolExist",
	}
	TopoStatusCode_value = map[string]int32{
		"Success":                            0,
		"InternalError":                      -1,
		"InvalidParam":                       -2,
		"InitFail":                           -3,
		"StorgeFail":                         -4,
		"IdDuplicated":                       -5,
		"ChunkServerNotFound":                -6,
		"ServerNotFound":                     -7,
		"ZoneNotFound":                       -8,
		"PhysicalPoolNotFound":               -9,
		"LogicalPoolNotFound":                -10,
		"CopySetNotFound":                    -11,
		"GenCopysetErr":                      -12,
		"AllocateIdFail":                     -13,
		"CannotRemoveWhenNotEmpty":           -14,
		"IpPortDuplicated":                   -15,
		"NameDuplicated":                     -16,
		"CreateCopysetNodeOnChunkServerFail": -17,
		"CannotRemoveNotRetired":             -18,
		"LogicalPoolExist":                   -19,
	}
)
