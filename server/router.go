package server

import (
	"os"

	"github.com/j75689/easybot/server/context"
	"github.com/j75689/easybot/server/middleware"

	rice "github.com/GeertJohan/go.rice"
	"github.com/foolin/gin-template/supports/gorice"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

func initRouter() (router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)

	router = gin.New()
	// middlerware
	router.HandleMethodNotAllowed = true
	router.NoMethod(middleware.NoMethodHandler())
	router.NoRoute(middleware.NoRouteHandler())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// auth
	router.Use(middleware.UserAuthMiddleware(middleware.AllowPathPrefixSkipper(
		middleware.Prefix{Method: "GET", Path: "/public"},
		middleware.Prefix{Method: "GET", Path: "/login"},
		middleware.Prefix{Method: "Any", Path: "/api/v1/login"},
		middleware.Prefix{Method: "Any", Path: "/api/v1/bot/hook"})))

	// Front-end file
	if _, err := os.Stat("dashboard/build"); err == nil {
		logger.Info("Serve dashboard/build")
		router.LoadHTMLGlob("./dashboard/build/index.html")
		router.Static("/public", "dashboard/build")
	} else {
		logger.Info("Serve rice")
		router.HTMLRender = gorice.New(rice.MustFindBox("../dashboard/build"))
		dist := rice.MustFindBox("../dashboard/build")
		router.StaticFS("/public", dist.HTTPBox())
	}

	// Register dashboard
	registerDashBoardRouter(router)

	// Register API
	if handler, err := httphandler.New(channel_secret, channel_token); err == nil {
		if bot, err := handler.NewClient(); err == nil {
			registerAPIRouter(router, handler, bot)
		} else {
			logger.Error(err)
		}
	} else {
		logger.Error(err)
	}

	return
}

func registerDashBoardRouter(app *gin.Engine) {
	dashboard := app.Group("/")
	dashboard.GET("/", context.HandleIndexPage(context_path))
	dashboard.GET("/dashboard", context.HandleIndexPage(context_path))
	dashboard.GET("/login", context.HandleIndexPage(context_path))
}

func registerAPIRouter(app *gin.Engine, handler *httphandler.WebhookHandler, botClient *linebot.Client) {
	v1 := app.Group("/api/v1")

	// login api
	v1.POST("/login", context.HandleLogin())

	// crud config
	v1.GET("/config/:id", context.HandleGetConfig(&db))
	v1.POST("/config/:id", context.HandlePostConfig(&db))
	v1.DELETE("/config/:id", context.HandleDeleteConfig(&db))

	// tester
	v1.POST("/runner", context.HandleTestRunner())
	v1.POST("/plugin/:plugin", context.HandleTestPlugin())

	botAPI := v1.Group("/bot")
	// Reply Hook
	botAPI.Any("/hook", context.HandleLineHook(handler, botClient))
	// Push Message
	botAPI.POST("/push/:userID", context.HandlePushMessage(botClient))
	botAPI.POST("/multicast", context.HandleMulticastMessage(botClient))

}
