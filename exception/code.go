package exception

type Code int

const (
	CodeParamsError Code = iota + 101
	CodeSystemError
	CodeBuildTokenError
	CodeTotpError
	CodeUAError
	CodeExistError
)
