package server

import (
	"os"

	rice "github.com/GeertJohan/go.rice"
	"github.com/foolin/gin-template/supports/gorice"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/j75689/easybot/server/context"
	"github.com/j75689/easybot/server/middleware"

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
	router.Use(middleware.TraceMiddleware(nil))

	// Register dashboard
	registerDashBoardRouter(router)

	// Register API
	var (
		lineHandler *httphandler.WebhookHandler
		lineBot     *linebot.Client
	)
	lineHandler, err := httphandler.New(channel_secret, channel_token)
	if err != nil {
		logger.Error("[Init] ", err)
	} else {
		if lineBot, err = lineHandler.NewClient(); err != nil {
			logger.Error("[Init] ", err)
		}
	}

	registerAPIRouter(router, lineHandler, lineBot)

	return
}

func registerDashBoardRouter(app *gin.Engine) {
	// static file
	if _, err := os.Stat("dashboard/build"); err == nil {
		logger.Info("[Init] ", "Serve dashboard/build")
		app.LoadHTMLGlob("./dashboard/build/index.html")
		app.Static("/public", "dashboard/build")
	} else {
		logger.Info("[Init] ", "Serve rice")
		app.HTMLRender = gorice.New(rice.MustFindBox("../dashboard/build"))
		dist := rice.MustFindBox("../dashboard/build")
		app.StaticFS("/public", dist.HTTPBox())
	}

	dashboard := app.Group("/")
	// session
	store := cookie.NewStore([]byte(appSecret))
	dashboard.Use(sessions.Sessions("session", store))
	// middleware
	dashboard.Use(middleware.SessionMiddleware(context_path))
	// page router
	dashboard.GET("/", context.HandleIndexPage(context_path))
	dashboard.GET("/dashboard", context.HandleIndexPage(context_path))
	dashboard.GET("/config", context.HandleIndexPage(context_path))
	dashboard.GET("/accessrole", context.HandleIndexPage(context_path))
	// login
	dashboard.GET("/login", context.HandleIndexPage(context_path))
	dashboard.POST("/login", context.HandleLogin(admin_user, admin_pass))
	// plugin
	dashboard.POST("/plugin/:plugin", context.HandleTestPlugin())
	// config
	dashboard.GET("/handler/config", context.HandleGetAllConfigID(&db))
	dashboard.GET("/handler/config/:id", context.HandleGetConfig(&db))
	dashboard.POST("/handler/config/:id", context.HandleCreateConfig(&db))
	dashboard.PUT("/handler/config/:id", context.HandleSaveConfig(&db))
	dashboard.DELETE("/handler/config/:id", context.HandleDeleteConfig(&db))
	// runner
	dashboard.POST("/handler/runner", context.HandleTestRunner())
	// accessrole
	dashboard.GET("/role/scope", context.HandleGetScopeTags())
	dashboard.GET("/role/account", context.HandleGetAllServiceAccount(&db))
	dashboard.DELETE("/role/account", context.HandleBatchDeleteServiceAccount(&db))
	dashboard.GET("/role/account/:name", context.HandleGetServiceAccount(&db))
	dashboard.PUT("/role/account/:name", context.HandleSaveServiceAccount(&db))
	dashboard.POST("/role/account/:name", context.HandleCreateServiceAccount(&db))
	dashboard.POST("/role/account/:name/refresh", context.HandleRefreshServiceAccountToken(&db))
}

func registerAPIRouter(app *gin.Engine, handler *httphandler.WebhookHandler, botClient *linebot.Client) {
	v1 := app.Group("/api/v1")
	// middleware
	skipper := middleware.AllowPathPrefixSkipper(
		middleware.Prefix{Method: "Any", Path: "/api/v1/bot/hook"})
	v1.Use(middleware.UserAuthMiddleware(skipper))
	v1.Use(middleware.UserIptableMiddleware(&db, skipper))
	v1.Use(middleware.UserScopeMiddleware(&db, skipper))

	// crud config
	v1.GET("/config/:id", context.HandleGetConfig(&db))
	v1.POST("/config/:id", context.HandleCreateConfig(&db))
	v1.PUT("/config/:id", context.HandleSaveConfig(&db))
	v1.DELETE("/config/:id", context.HandleDeleteConfig(&db))

	// plugin
	v1.POST("/plugin/:plugin", context.HandleTestPlugin())

	botAPI := v1.Group("/bot")
	// Reply Hook
	botAPI.Any("/hook", context.HandleLineHook(handler, botClient))
	// Push Message
	botAPI.POST("/push/:userID", context.HandlePushMessage(botClient))
	botAPI.POST("/multicast", context.HandleMulticastMessage(botClient))

}
