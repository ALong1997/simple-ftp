package common

import "errors"

// Error
var (
	NoErr             = errors.New("OK")                             // OK
	InnerErr          = errors.New("inner err")                      // 内部错误
	InvalidCommandErr = errors.New("command not found")              // 无效命令
	AuthenticationErr = errors.New("authentication failure")         // 认证失败
	PutDirErr         = errors.New("put 命令不支持发送文件夹，请尝试putdir命令")     // put 文件夹
	CDArgsErr         = errors.New("cd parameter must be directory") // cd 参数必须是文件目录
)
