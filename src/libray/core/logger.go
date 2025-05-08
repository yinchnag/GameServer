// JoysGames copyrights this specification. No part of this specification may be
// reproduced in any form or means, without the prior written consent of JoysGames.
//
// This specification is preliminary and is subject to change at any time without notice.
// JoysGames assumes no responsibility for any errors contained herein.
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
// @package JGServer
// @copyright joysgames.cn All rights reserved.
// @version v1.0

package core

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	Logger  *logrus.Logger      = defLogger()
	loggers map[string]*LoggerS = make(map[string]*LoggerS, 0) // 日志列表
	lock    sync.RWMutex                                       // 读写锁
)

// 默认日志器
func defLogger() *logrus.Logger {
	return GetLogger("log", func() *JS_LoggerConfig { return nil })
}

// 设置默认日志器
func SetDefLogger(logger *logrus.Logger) {
	if logger != nil {
		Logger = logger
	}
}

// 获取日志器
func GetLogger(name string, loader func() *JS_LoggerConfig) *logrus.Logger {
	lock.RLock()
	logger, ok := loggers[name]
	lock.RUnlock()
	if !ok {
		lock.Lock()
		logger = new(LoggerS)
		logger.Name = name
		logger.loader = loader
		logger.logger = logrus.New()
		logger.logger.Formatter = &TextFormatter{} // 设置默认格式化
		logger.initLogger()
		loggers[name] = logger
		lock.Unlock()
	}
	return logger.logger
}

// 日志接口(多线程安全)
type LoggerS struct {
	Name   string                  // 日志名称
	loader func() *JS_LoggerConfig // 日志配置名称
	logger *logrus.Logger          // 日志器皿
}

// 初始化日志器
func (that *LoggerS) initLogger() {
	conf := that.loader()
	if conf == nil {
		return
	}
	// if !IsDebug {
	// 	conf.LogLevel = HF_MinInt(4, conf.LogLevel) // 非调试模式强制关闭debug信息
	// }
	that.logger.SetLevel(logrus.Level(conf.LogLevel))

	// 普通日志
	logName := fmt.Sprintf("%s/logs/%s.log", GetBasePath(), that.Name)
	logWriter, err := rotatelogs.New(
		logName+".%Y%m%d",                         // 分割后的文件名称
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 设置最大保存时间(7天)
		rotatelogs.WithRotationTime(24*time.Hour), // 设置日志切割时间间隔(1天)
	)
	if err != nil {
		log.Println("failed to create new rotatelogs")
		return
	}

	// 错误日志
	logName = fmt.Sprintf("%s/logs/%s-error.log", GetBasePath(), that.Name)
	errorWriter, err := rotatelogs.New(
		logName+".%Y%m%d",                         // 分割后的文件名称
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 设置最大保存时间(7天)
		rotatelogs.WithRotationTime(24*time.Hour), // 设置日志切割时间间隔(1天)
	)
	if err != nil {
		log.Println("failed to create new rotatelogs")
		return
	}

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: errorWriter,
		logrus.PanicLevel: errorWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &TextFormatter{})
	that.logger.SetReportCaller(conf.LogConsole)
	that.logger.AddHook(lfHook)
}

// 日志格式化
type TextFormatter struct{}

// 日志格式化
func (that *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buff *bytes.Buffer
	if entry.Buffer != nil {
		buff = entry.Buffer
	} else {
		buff = &bytes.Buffer{}
	}

	// 格式化日志信息
	var newlog string
	timestamp := ServerTime().Format(DATE_FORMAT1)
	level := strings.ToUpper(entry.Level.String())
	goroutineID := fmt.Sprintf("00000000%d", GetGoroutineID())
	if entry.HasCaller() {
		fname := filepath.Base(entry.Caller.File)
		newlog = fmt.Sprintf("[:%s] %s %s:%d [%s] %s\n", goroutineID[len(goroutineID)-8:], timestamp, fname, entry.Caller.Line, level, entry.Message)
	} else {
		newlog = fmt.Sprintf("[:%s] %s [%s] %s\n", goroutineID[len(goroutineID)-8:], timestamp, level, entry.Message)
	}
	buff.WriteString(newlog)
	return buff.Bytes(), nil
}
