package server

import (
	"github.com/j75689/easybot/config"
	"github.com/j75689/easybot/pkg/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/foolin/gin-template/supports/gorice"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"go.uber.org/zap"

	bolt "go.etcd.io/bbolt"
)

func initRouter(isDebug *bool) (router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)

	router = gin.New()
	//router.Use(ginzap.Ginzap(defaultlogger, time.RFC3339, true))
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	//router.LoadHTMLGlob("./dashboard/build/index.html")
	router.HTMLRender = gorice.New(rice.MustFindBox("../dashboard/build"))
	// Session
	store := cookie.NewStore([]byte("easybot"))
	router.Use(sessions.Sessions("session", store))

	// LineBot Hook
	router.Any(webhook_path, gin.WrapH(newLineHookHandler()))

	// Front-end file
	if _, err := os.Stat("dashboard/build"); err == nil {
		logger.Info("Serve dashboard/build")
		router.Static("/public", "dashboard/build")
	} else {
		logger.Info("Serve rice ", err)
		dist := rice.MustFindBox("../dashboard/build")
		router.StaticFS("/public", dist.HTTPBox())
	}

	ManagerGroup := router.Group("/")
	ManagerGroup.Use(SessionMiddleware())
	// Dashboard
	{
		// API
		ManagerGroup.POST("/login", handleLogin)
		// Pages
		ManagerGroup.GET("/", handleIndexPage)
		ManagerGroup.GET("/dashboard", handleIndexPage)
		ManagerGroup.GET("/login", handleIndexPage)
	}

	// Config CRUD
	router.Any("/config/:id", handleCRUDConfig)

	// Tester
	if *isDebug {
		router.Any("/debug/runner", handleTestRunner)
		router.Any("/debug/plugin/:plugin", handleTestPlugin)
		router.Any("/debug/request", handleTestRequest)
	}

	return
}

func handleIndexPage(c *gin.Context) {

	if c.Request.URL.Path == "/" {
		c.Redirect(301, context_path+"dashboard")
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func newLineHookHandler() (handler *httphandler.WebhookHandler) {
	// Line SDK
	handler, err := httphandler.New(
		channel_secret,
		channel_token,
	)
	if err != nil {
		logger.Error(err)
		return
	}

	// Setup HTTP Server for receiving requests from LINE platform
	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {

		logger.Info(r)

		bot, err := handler.NewClient()
		if err != nil {
			logger.Error(err)
		}
		for _, event := range events {
			if msg := handleMessage(event.Source, event.Message); msg != nil {
				msgData, _ := msg.MarshalJSON()
				logger.Debug(string(msgData))
				if _, err = bot.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
					logger.Error(err)
				}
			}
		}
	})

	return
}

func handleCRUDConfig(c *gin.Context) {
	var (
		bucketName = "easybot.config"
		configID   = c.Param("id")
	)

	switch c.Request.Method {
	case "GET":

		db.Batch(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			if err != nil {
				logger.Errorf("open bucket: %s", err)
			}
			v := b.Get([]byte(configID))
			if v != nil {
				var config map[string]interface{}
				json.Unmarshal(v, &config)
				c.JSON(200, config)
			} else {
				c.JSON(200, map[string]string{"error": "Config Not Found"})
			}

			return nil
		})
	case "POST":
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			if err != nil {
				logger.Errorf("open bucket: %s", err)
				return err
			}
			if configData, err := c.GetRawData(); err == nil {
				err = registerMessageHandlerConfig(configID, configData)
				if err != nil {
					logger.Errorf("register config [%s]: %s", configID, err)
					c.JSON(200, map[string]string{"error": err.Error()})
					return err
				}
				err = b.Put([]byte(configID), configData)
				if err != nil {
					logger.Errorf("save config [%s]: %s", configID, err)
					return err
				}

				logger.Infof("save config [%s] success.", configID)
			}

			return nil
		})
	case "DELETE":
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			if err != nil {
				logger.Errorf("open bucket: %s", err)
				return err
			}

			v := b.Get([]byte(configID))
			if v != nil {
				var messageConfig config.MessageHandlerConfig
				err = json.Unmarshal(v, &messageConfig)
				if err == nil {
					unregisterMessageHandlerConfig(&messageConfig)
				}
			}

			err = b.Delete([]byte(configID))
			if err != nil {
				logger.Errorf("delete config [%s] failed.", configID)
				return err
			}
			logger.Infof("delete config [%s] success.", configID)
			return nil
		})
	default:
		c.JSON(405, map[string]string{"error": "Method Not Allowed."})
	}
}

func handleTestRunner(c *gin.Context) {
	type args struct {
		Message   string                 `json:"message"`
		Variables map[string]interface{} `json:"variables"`
	}
	defer c.Done()
	postdata, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	var arg args
	err = json.Unmarshal(postdata, &arg)
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(arg)
	reply := handleTextMessage(arg.Message, &arg.Variables)

	c.JSON(200, reply)
}

func handleTestPlugin(c *gin.Context) {
	var (
		pluginName = c.Param("plugin")
	)
	defer c.Done()
	if f, ok := pluginfuncs.Load(pluginName); ok {
		plugin := f.(*config.PluginFunc)

		type args struct {
			Input     interface{}            `json:"input"`
			Variables map[string]interface{} `json:"variables"`
			logger    *zap.SugaredLogger
		}

		postdata, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error(err)
			return
		}

		var arg = args{
			logger: logger,
		}
		err = json.Unmarshal(postdata, &arg)
		if err != nil {
			logger.Error(err)
		}
		logger.Debug(arg)
		var b []byte
		if b, err = json.Marshal(arg.Input); err != nil {
			logger.Error(err.Error())
		}
		ParamData := util.ReplaceVariables(string(b), arg.Variables)
		if err = json.Unmarshal([]byte(ParamData), &arg.Input); err != nil {
			logger.Error(err.Error())
		}
		v, err := (*plugin)(arg.Input, arg.Variables, arg.logger)
		if err != nil {
			c.JSON(200, map[string]interface{}{"error": err.Error()})
		} else {
			c.JSON(200, v)
		}

	}

}

func handleTestRequest(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	logger.Info(string(body))
	c.JSON(200, map[string]interface{}{
		"headers": c.Request.Header,
		"body":    string(body),
	})
}

func handleLogin(c *gin.Context) {
	var (
		user   = c.DefaultPostForm("user", "")
		pass   = c.DefaultPostForm("pass", "")
		result = map[string]interface{}{
			"success": false,
		}
	)

	if user == admin_user && pass == admin_pass {
		session := sessions.Default(c)
		session.Set("login", true)
		session.Save()
		result["success"] = true
	}

	c.JSON(200, result)
}

func SessionMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		var isLogin = false
		session := sessions.Default(c)
		v := session.Get("login")
		if v != nil {
			isLogin = v.(bool)
		}
		if isLogin || strings.Index(c.Request.URL.Path, "/login") > -1 {
			c.Next()
		} else {
			c.Redirect(http.StatusMovedPermanently, context_path+"login")
			c.Done()
		}

	}

}
