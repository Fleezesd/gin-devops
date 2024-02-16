package common

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewZapLogger(logLevel string, logFilePath string) *zap.Logger {

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch logLevel {
	case "DEBUG":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "INFO":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	default:
		atomicLevel.SetLevel(zapcore.InfoLevel)
	}

	// 日志轮转
	writer := &lumberjack.Logger{
		Filename:   logFilePath, // 日志名称
		MaxSize:    1,           // 日志大小限制，单位MB
		MaxAge:     30,          // 历史日志文件保留天数
		MaxBackups: 2,           // 最大保留历史日志数量
		LocalTime:  true,        // 本地时区
		Compress:   false,       // 历史日志文件压缩标识
	}

	// 设置默认字段
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "line",                                                 // 调用者
		MessageKey:     "msg",                                                  // 内容msg
		FunctionKey:    "func",                                                 // 函数func
		StacktraceKey:  "stacktrace",                                           // 堆栈stackTrace
		LineEnding:     zapcore.DefaultLineEnding,                              // 换行字符 \n
		EncodeLevel:    zapcore.LowercaseLevelEncoder,                          // 小写字符
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"), // 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 打印控制台用的Core
	zapCoreConsole := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 日志字段配置
		zapcore.AddSync(os.Stdout),            // 控制台输出
		atomicLevel,                           // 日志级别
	)

	// 写文件用的CoreFile
	zapCoreFile := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		zap.InfoLevel,
	)

	// 底层核心 分配给不同通道(控制台/文件)
	core := zapcore.NewTee(
		zapCoreConsole,
		zapCoreFile,
	)

	return zap.New(core, zap.AddCaller()) // 启用调用者信息功能
}
