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

type JS_DatabaseConf struct {
	DBUser      string `json:"dbuser"`      // 游戏数据库
	DBLog       string `json:"dblog"`       // 日志数据库
	Redis       string `json:"redis"`       // redis地址
	RedisDB     int    `json:"redisdb"`     // redis db编号
	RedisAuth   string `json:"redisauth"`   // redis认证
	RedisPrefix string `json:"redisprefix"` // redis前缀
	MaxDBConn   int    `json:"maxdbconn"`   // 最大数据库连接数
}
