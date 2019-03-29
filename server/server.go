package server

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/j75689/easybot/handler"
	"github.com/j75689/easybot/plugin"
	"go.uber.org/zap"

	"github.com/j75689/easybot/config"

	log "github.com/j75689/easybot/pkg/logger"
	"github.com/j75689/easybot/pkg/store"
)

var (
	appName        = "easybot"
	channel_secret = os.Getenv("CHANNEL_SECRET")
	channel_token  = os.Getenv("CHANNEL_TOKEN")
	port           = os.Getenv("PORT")
	webhook_path   = os.Getenv("WEBHOOK_PATH")
	plugin_path    = os.Getenv("PLUGIN_PATH")
	db_driver      = os.Getenv("DB_DRIVER")
	db_path        = os.Getenv("DB_PATH")
	context_path   = os.Getenv("CONTEXT_PATH")
	admin_user     = os.Getenv("ADMIN_USER")
	admin_pass     = os.Getenv("ADMIN_PASS")
	loggerLevel    = os.Getenv("LOG_LEVEL")
	loggerPath     = os.Getenv("LOG_PATH")

	logger *zap.SugaredLogger
	db     store.Storage
)

func initServer() {
	if port == "" {
		port = "8801"
	}
	if webhook_path == "" {
		webhook_path = "/webhook"
	}
	if plugin_path == "" {
		plugin_path = "./plugin"
	}
	if db_path == "" {
		db_path = "./data/" + appName + ".db"
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

	logger.Info("init db")
	// init db
	var err error
	db, err = store.NewStoreage(db_driver, &store.Connection{
		DBName: "config",
		Host:   db_path,
	})
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("init config")
	// init config
	if err = db.LoadAll(func(key string, value interface{}) {
		if b, err := json.Marshal(value); err == nil {
			var cfg config.MessageHandlerConfig
			if err = json.Unmarshal(b, &cfg); err == nil {
				handler.RegisterConfig(&cfg)
			} else {
				logger.Errorf("Unmarshal config [%s] error: %v", key, err)
			}
		}
	}); err != nil {
		logger.Error(err)
	}
}

// Start 啟動服務
func Start(mode string) {
	initServer()
	logger.Infof("Service start on localhost:%s", port)
	if err := initRouter().Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.Error(err)
	}

}
