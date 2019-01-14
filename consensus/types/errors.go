package types

var (
	OK = NewResultOK(nil, "")

	ErrInternalError     = NewError(CodeType_InternalError, "Internal error")
	ErrEncodingError     = NewError(CodeType_EncodingError, "Encoding error")
	ErrBadNonce          = NewError(CodeType_BadNonce, "Error bad nonce")
	ErrUnauthorized      = NewError(CodeType_Unauthorized, "Unauthorized")
	ErrInsufficientFunds = NewError(CodeType_InsufficientFunds, "Insufficient funds")
	ErrUnknownRequest    = NewError(CodeType_UnknownRequest, "Unknown request")

	ErrBaseInsufficientFees     = NewError(CodeType_BaseInsufficientFees, "Error (base) insufficient fees")
	ErrBaseInvalidInput         = NewError(CodeType_BaseInvalidInput, "Error (base) invalid input")
	ErrBaseInvalidSignature     = NewError(CodeType_BaseInvalidSignature, "Error (base) invalid signature")
	ErrBaseUnknownAddress       = NewError(CodeType_BaseUnknownAddress, "Error (base) unknown address")
)
