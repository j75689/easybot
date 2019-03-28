package server

import (
	"github.com/j75689/easybot/config"
	"github.com/j75689/easybot/plugin"
	"fmt"
	"os"
	"strconv"
	"sync"

	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	defaultlogger  *zap.Logger
	logger         *zap.SugaredLogger
	channel_secret = os.Getenv("CHANNEL_SECRET")
	channel_token  = os.Getenv("CHANNEL_TOKEN")
	port           = 8801
	webhook_path   = "/webhook"
	plugin_path    = "./plugin"
	db_path        = "./data/easybot.db"
	context_path   = os.Getenv("CONTEXT_PATH")
	admin_user     = "admin"
	admin_pass     = "admin"

	db          *bolt.DB
	pluginfuncs = &sync.Map{}
	// TextMessageHandleConfig 文字訊息的處理設定檔
	TextMessageHandleConfig = &sync.Map{}
)

func initService(isDebug *bool) {
	// setting logger
	{
		hook := lumberjack.Logger{
			Filename:   "./logs/easybot.log", // 日誌文件路徑
			MaxSize:    128,                     // 每個日誌文件保存的最大尺寸 單位：M
			MaxBackups: 30,                      // 日誌文件最多保存多少個備份
			MaxAge:     7,                       // 文件最多保存多少天
			Compress:   true,                    // 是否壓縮
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
		if *isDebug {
			atomicLevel.SetLevel(zap.DebugLevel)
		} else {
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
	logger.Info("init env")
	if os.Getenv("PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("PORT"))
	}
	if os.Getenv("WEBHOOK_PATH") != "" {
		webhook_path = os.Getenv("WEBHOOK_PATH")
	}
	if os.Getenv("PLUGIN_PATH") != "" {
		plugin_path = os.Getenv("PLUGIN_PATH")
	}
	if os.Getenv("DB_PATH") != "" {
		db_path = os.Getenv("DB_PATH")
	}
	if os.Getenv("ADMIN_USER") != "" {
		admin_user = os.Getenv("ADMIN_USER")
	}
	if os.Getenv("ADMIN_PASS") != "" {
		admin_pass = os.Getenv("ADMIN_PASS")
	}
	logger.Info("load plugin")
	// add default plugin
	{
		graphql := config.PluginFunc(plugin.Graphql)
		pluginfuncs.Store("Graphql", &graphql)
		equal := config.PluginFunc(plugin.Equal)
		pluginfuncs.Store("Equal", &equal)
	}
	// load addition plugin
	LoadPlugins(plugin_path)

	logger.Info("init db")
	// init db
	var err error
	db, err = bolt.Open(db_path, 0644, nil)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("init config")
	// init config
	db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("easybot.config"))
		if err != nil {
			logger.Errorf("open bucket: %s", err)
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := registerMessageHandlerConfig(string(k), v); err != nil {
				logger.Error(err)
			}

		}

		return nil
	})
}

// Start 啟動服務
func Start(mode string) {
	var (
		isDebug = mode != "production"
	)
	initService(&isDebug)

	logger.Infof("Service start on localhost:%d", port)
	if err := initRouter(&isDebug).Run(fmt.Sprintf(":%d", port)); err != nil {
		logger.Error(err)
	}

}
