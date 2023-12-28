package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zhfreal/E5SubBot/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志切割设置
// func getLogWriter(log_dir string) zapcore.WriteSyncer {
// 	lumberJackLogger := &lumberjack.Logger{
// 		Filename:   filepath.Join(log_dir, "log/latest.log"), // 日志文件位置
// 		MaxSize:    1,                                        // 日志文件最大大小(MB)
// 		MaxBackups: 5,                                        // 保留旧文件最大数量
// 		MaxAge:     30,                                       // 保留旧文件最长天数
// 		Compress:   true,                                     // 是否压缩旧文件
// 	}
// 	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
// }

// generate zapcore.WriteSyncer
func getLogWriter(log_into_file bool, log_file string, max_size, max_backups, max_age int) zapcore.WriteSyncer {
	if log_into_file {
		// log logrotate settings
		lumberJackLogger := &lumberjack.Logger{
			Filename:   log_file,    // 日志文件位置
			MaxSize:    max_size,    // 日志文件最大大小(MB)
			MaxBackups: max_backups, // 保留旧文件最大数量
			MaxAge:     max_age,     // 保留旧文件最长天数
			Compress:   true,        // 是否压缩旧文件
			LocalTime:  true,        // 是否使用本地时间
		}
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(os.Stdout)
}

// 编码器
func getEncoder() zapcore.Encoder {
	// 使用默认的JSON编码
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// Init 初始化Logger
// func Init(work_dir string) {
// 	writeSyncer := getLogWriter(work_dir)
// 	encoder := getEncoder()
// 	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
// 	zap.ReplaceGlobals(zap.New(core, zap.AddCaller()))
// }

// we support debug, info, warn, error logs
func Init(log_into_file bool, log_dir, log_file, log_level string, max_size, max_backups, max_age int) {
	var t_log_path string
	if log_into_file {
		t_log_path = filepath.Join(log_dir, log_file)
		err := utils.MakeFolderAndSetPermission(t_log_path)
		if err != nil {
			fmt.Printf("failed to make folder and set it's permission to 755, failed with: %v\n", err.Error())
			panic(err)
		}
	} else {
		t_log_path = ""
	}
	writeSyncer := getLogWriter(log_into_file, t_log_path, max_size, max_backups, max_age)
	encoder := getEncoder()
	// convert log_level into lower case
	log_level = strings.ToLower(log_level)
	var zap_log_level zapcore.Level
	switch log_level {
	case "debug":
		zap_log_level = zapcore.DebugLevel
	case "info":
		zap_log_level = zapcore.InfoLevel
	case "warn":
		zap_log_level = zapcore.WarnLevel
	case "error":
		zap_log_level = zapcore.ErrorLevel
	default:
		zap_log_level = zapcore.WarnLevel
	}
	core := zapcore.NewCore(encoder, writeSyncer, zap_log_level)
	zap.ReplaceGlobals(zap.New(core, zap.AddCaller()))
}

func Debug(msg ...interface{}) {
	zap.S().Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	zap.S().Debugf(format, args...)
}

func Debugln(msg ...interface{}) {
	zap.S().Debugln(msg)
}

func Info(msg ...interface{}) {
	zap.S().Info(msg)
}

func Infof(format string, args ...interface{}) {
	zap.S().Infof(format, args...)
}

func Infoln(msg ...interface{}) {
	zap.S().Infoln(msg)
}

func Warn(msg ...interface{}) {
	zap.S().Warn(msg)
}

func Warnf(format string, args ...interface{}) {
	zap.S().Warnf(format, args...)
}

func Warnln(msg ...interface{}) {
	zap.S().Warnln(msg)
}

func Error(msg ...interface{}) {
	zap.S().Error(msg)
}

func Errorf(format string, args ...interface{}) {
	zap.S().Errorf(format, args...)
}

func Errorln(msg ...interface{}) {
	zap.S().Errorln(msg)
}

func Fatal(msg ...interface{}) {
	zap.S().Fatal(msg)
}

func Fatalf(format string, args ...interface{}) {
	zap.S().Fatalf(format, args...)
}

func Fatalln(msg ...interface{}) {
	zap.S().Fatalln(msg)
}
