package core

import "os"

var (
	IsDebug   = os.Getenv("GIN_MODE") != ""
	IsTesting = false
)

type JS_LoggerConfig struct {
	MaxFileSize int64 `json:"maxfilesize"` // 文件长度
	MaxFileNum  int   `json:"maxfilenum"`  // 文件数量
	LogLevel    int   `json:"loglevel"`    // 日志等级
	LogConsole  bool  `json:"logconsole"`  // 是否输出控制台
	LogCaller   bool  `json:"logcaller"`   // 是否输出文件源
}
