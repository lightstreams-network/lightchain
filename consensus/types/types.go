package types

import "fmt"
import (
	"math"
	"strconv"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = fmt.Errorf
var _ = math.Inf


// Not being used
// Could be added to request/response
// so we don't have to type switch
// (would be twice as fast, but we're talking about 15ns)
type MessageType int32

var MessageType_name = map[int32]string{
	0:  "NullMessage",
	1:  "Echo",
	2:  "Flush",
	3:  "Info",
	4:  "SetOption",
	5:  "Exception",
	17: "DeliverTx",
	18: "CheckTx",
	19: "Commit",
	20: "Query",
	21: "InitChain",
	22: "BeginBlock",
	23: "EndBlock",
}

func (x MessageType) String() string {
	return EnumName(MessageType_name, int32(x))
}

func EnumName(m map[int32]string, v int32) string {
	s, ok := m[v]
	if ok {
		return s
	}
	return strconv.Itoa(int(v))
}

type CodeType int32

const (
	CodeType_OK CodeType = 0
	// General response codes, 0 ~ 99
	CodeType_InternalError     CodeType = 1
	CodeType_EncodingError     CodeType = 2
	CodeType_BadNonce          CodeType = 3
	CodeType_Unauthorized      CodeType = 4
	CodeType_InsufficientFunds CodeType = 5
	CodeType_UnknownRequest    CodeType = 6
	// Reserved for basecoin, 100 ~ 199
	CodeType_BaseDuplicateAddress     CodeType = 101
	CodeType_BaseEncodingError        CodeType = 102
	CodeType_BaseInsufficientFees     CodeType = 103
	CodeType_BaseInsufficientFunds    CodeType = 104
	CodeType_BaseInsufficientGasPrice CodeType = 105
	CodeType_BaseInvalidInput         CodeType = 106
	CodeType_BaseInvalidOutput        CodeType = 107
	CodeType_BaseInvalidPubKey        CodeType = 108
	CodeType_BaseInvalidSequence      CodeType = 109
	CodeType_BaseInvalidSignature     CodeType = 110
	CodeType_BaseUnknownAddress       CodeType = 111
	CodeType_BaseUnknownPubKey        CodeType = 112
	CodeType_BaseUnknownPlugin        CodeType = 113
)

