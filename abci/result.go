package abci

func NewError(code CodeType, log string) Result {
	return Result{
		Code: code,
		Log:  log,
	}
}

// CONTRACT: a zero Result is OK.
type Result struct {
	Code CodeType `json:"code"`
	Log  string   `json:"log"` // Can be non-deterministic
}

// NOTE: if data == nil and log == "", same as zero Result.
func NewResultOK(data []byte, log string) Result {
	return Result{
		Code: CodeType_OK,
		Log:  log,
	}
}
