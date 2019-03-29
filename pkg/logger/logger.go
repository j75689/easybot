package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	defaultlogger *zap.Logger
	logger        *zap.SugaredLogger
)

// NewLogger instance
func NewLogger(appName, loggerPath, loggerLevel string) {

	// setting logger
	{
		hook := lumberjack.Logger{
			Filename:   loggerPath + appName + ".log", // 日誌文件路徑
			MaxSize:    128,                           // 每個日誌文件保存的最大尺寸 單位：M
			MaxBackups: 30,                            // 日誌文件最多保存多少個備份
			MaxAge:     7,                             // 文件最多保存多少天
			Compress:   true,                          // 是否壓縮
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "line",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 彩色編碼器
			EncodeTime:     zapcore.ISO8601TimeEncoder,       // ISO8601 UTC 時間格式
			EncodeDuration: zapcore.SecondsDurationEncoder,   //
			EncodeCaller:   zapcore.ShortCallerEncoder,       // 路徑編碼器
			EncodeName:     zapcore.FullNameEncoder,
		}

		// 設置日誌級別
		atomicLevel := zap.NewAtomicLevel()
		switch strings.ToUpper(loggerLevel) {
		case "DEBUG":
			atomicLevel.SetLevel(zap.DebugLevel)
		case "INFO":
			atomicLevel.SetLevel(zap.InfoLevel)
		case "WARN":
			atomicLevel.SetLevel(zap.WarnLevel)
		case "ERROR":
			atomicLevel.SetLevel(zap.ErrorLevel)
		case "DPANIC":
			atomicLevel.SetLevel(zap.DPanicLevel)
		case "PANIC":
			atomicLevel.SetLevel(zap.PanicLevel)
		case "FATAL":
			atomicLevel.SetLevel(zap.FatalLevel)
		default:
			atomicLevel.SetLevel(zap.InfoLevel)
		}

		core := zapcore.NewCore(
			//zapcore.NewJSONEncoder(encoderConfig), // 編碼器配置
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制枱和文件
			atomicLevel, // 日誌級別
		)

		// 開啟開發模式，堆棧跟蹤
		caller := zap.AddCaller()
		// 開啟文件及行號
		development := zap.Development()
		// 構造日誌
		defaultlogger := zap.New(core, development, caller)
		logger = defaultlogger.Sugar()
	}
}

// GetLogger get logger instance
func GetLogger() *zap.SugaredLogger {
	return logger
}

// Debug log debug level.
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf log debug level with template.
func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

// Debugw log debug level with custom field.
func Debugw(message string, keyValue ...interface{}) {
	logger.Debugw(message, keyValue...)
}

// Info log info level.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof log info level with template.
func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

// Infow log info level with custom field.
func Infow(message string, keyValue ...interface{}) {
	logger.Infow(message, keyValue...)
}

// Warn log warn level.
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warnf log warn level with template.
func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

// Warnw log warn level with custom field.
func Warnw(message string, keyValue ...interface{}) {
	logger.Warnw(message, keyValue...)
}

// Error log error level.
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf log error level with template.
func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

// Errorw log error level with custom field.
func Errorw(message string, keyValue ...interface{}) {
	logger.Errorw(message, keyValue...)
}

// DPanic log dpanic level.
func DPanic(args ...interface{}) {
	logger.DPanic(args...)
}

// DPanicf log dpanic level with template.
func DPanicf(template string, args ...interface{}) {
	logger.DPanicf(template, args...)
}

// DPanicw log dpanic level with custom field.
func DPanicw(message string, keyValue ...interface{}) {
	logger.DPanicw(message, keyValue...)
}

// Panic log panic level.
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Panicf log panic level with template.
func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

// Panicw log panic level with custom field.
func Panicw(message string, keyValue ...interface{}) {
	logger.Panicw(message, keyValue...)
}

// Fatal log fatal level.
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Fatalf log fatal level with template.
func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

// Fatalw log fatal level with custom field.
func Fatalw(message string, keyValue ...interface{}) {
	logger.Fatalw(message, keyValue...)
}
