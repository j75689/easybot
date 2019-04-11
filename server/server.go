package server

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/j75689/easybot/server/auth"

	"github.com/j75689/easybot/handler"
	"github.com/j75689/easybot/plugin"
	"go.uber.org/zap"

	"github.com/j75689/easybot/config"

	log "github.com/j75689/easybot/pkg/logger"
	"github.com/j75689/easybot/pkg/store"
)

var (
	appName        = "easybot"
	appSecret      = os.Getenv("APP_SECRET")
	channel_secret = os.Getenv("CHANNEL_SECRET")
	channel_token  = os.Getenv("CHANNEL_TOKEN")
	port           = os.Getenv("PORT")
	plugin_path    = os.Getenv("PLUGIN_PATH")
	db_driver      = os.Getenv("DB_DRIVER")
	db_name        = os.Getenv("DB_NAME")
	db_host        = os.Getenv("DB_HOST")
	db_port        = os.Getenv("DB_PORT")
	db_user        = os.Getenv("DB_USER")
	db_pass        = os.Getenv("DB_PASS")
	context_path   = os.Getenv("CONTEXT_PATH")
	admin_user     = os.Getenv("ADMIN_USER")
	admin_pass     = os.Getenv("ADMIN_PASS")
	loggerLevel    = os.Getenv("LOG_LEVEL")
	loggerPath     = os.Getenv("LOG_PATH")

	logger *zap.SugaredLogger
	db     store.Storage
)

func initServer() {
	if appSecret == "" {
		appSecret = appName
	}
	if context_path == "" {
		context_path = "/"
	}
	if port == "" {
		port = "8801"
	}
	if plugin_path == "" {
		plugin_path = "./plugin"
	}
	if db_name == "" {
		db_name = appName
	}
	if db_host == "" {
		db_host = "./data/" + appName + ".db"
	}
	if db_driver == "" {
		db_driver = "bolt"
	}
	if admin_user == "" {
		admin_user = "admin"
	}
	if admin_pass == "" {
		admin_pass = "admin"
	}
	if loggerPath == "" {
		loggerPath = "./logs/"
	}
	// init logger
	logger = log.NewLogger(appName, loggerPath, loggerLevel)
	// init plugin
	plugin.Load(plugin_path, log.GetLogger())

	logger.Infof("init db driver:[%s]", db_driver)
	// init db
	var err error
	db, err = store.NewStoreage(db_driver, &store.Connection{
		DBName: db_name,
		Host:   db_host,
		Port:   db_port,
		User:   db_user,
		Pass:   db_pass,
	})
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("init config")
	// init config
	if err = db.LoadAll("config", func(key string, value interface{}) {
		if b, err := json.Marshal(value); err == nil {
			var cfg config.MessageHandlerConfig
			if err = json.Unmarshal(b, &cfg); err == nil {
				handler.RegisterConfig(&cfg)
				logger.Infof("Register config [%s]", key)
			} else {
				logger.Errorf("Unmarshal config [%s] error: %v", key, err)
			}
		}
	}); err != nil {
		logger.Error(err)
	}
	logger.Info("init Auth module")
	// init Auth module
	auth.SetSigningKey(appSecret)
}

// Start 啟動服務
func Start(mode string) {
	initServer()
	logger.Infof("Service start on localhost:%s", port)
	if err := initRouter().Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.Error(err)
	}

}
