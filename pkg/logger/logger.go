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

	Debug   func(args ...interface{})
	Debugf  func(template string, args ...interface{})
	Debugw  func(msg string, keysAndValues ...interface{})
	Info    func(args ...interface{})
	Infow   func(msg string, keysAndValues ...interface{})
	Infof   func(template string, args ...interface{})
	Error   func(args ...interface{})
	Errorf  func(template string, args ...interface{})
	Errorw  func(msg string, keysAndValues ...interface{})
	Warn    func(args ...interface{})
	Warnf   func(template string, args ...interface{})
	Warnw   func(msg string, keysAndValues ...interface{})
	Panic   func(args ...interface{})
	Panicf  func(template string, args ...interface{})
	Panicw  func(msg string, keysAndValues ...interface{})
	DPanic  func(args ...interface{})
	DPanicf func(template string, args ...interface{})
	DPanicw func(msg string, keysAndValues ...interface{})
	Fatal   func(args ...interface{})
	Fatalf  func(template string, args ...interface{})
	Fatalw  func(msg string, keysAndValues ...interface{})
)

// NewLogger instance
func NewLogger(appName, loggerPath, loggerLevel string) *zap.SugaredLogger {

	// setting logger
	{
		hook := lumberjack.Logger{
			Filename:   loggerPath + appName + ".log",
			MaxSize:    128,
			MaxBackups: 30,
			MaxAge:     7,
			Compress:   true,
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "line",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		}

		// Setting Log Level
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
			//zapcore.NewJSONEncoder(encoderConfig), // Json
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // console
			atomicLevel,
		)

		// Trace code
		caller := zap.AddCaller()
		// open code line
		development := zap.Development()
		// construct
		defaultlogger := zap.New(core, development, caller)
		logger = defaultlogger.Sugar()

		// public func
		Debug = logger.Debug
		Debugf = logger.Debugf
		Debugw = logger.Debugw
		Info = logger.Info
		Infow = logger.Infow
		Infof = logger.Infof
		Error = logger.Error
		Errorf = logger.Errorf
		Errorw = logger.Errorw
		Warn = logger.Warn
		Warnf = logger.Warnf
		Warnw = logger.Warnw
		Panic = logger.Panic
		Panicf = logger.Panicf
		Panicw = logger.Panicw
		DPanic = logger.DPanic
		DPanicf = logger.DPanicf
		DPanicw = logger.DPanicw
		Fatal = logger.Fatal
		Fatalf = logger.Fatalf
		Fatalw = logger.Fatalw
	}

	return logger
}

// GetLogger get logger instance
func GetLogger() *zap.SugaredLogger {
	return logger
}

// GetDefaultLogger get defaultlogger instance
func GetDefaultLogger() *zap.Logger {
	return defaultlogger
}
