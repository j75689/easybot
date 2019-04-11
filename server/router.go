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
	// static file
	if _, err := os.Stat("dashboard/build"); err == nil {
		logger.Info("Serve dashboard/build")
		app.LoadHTMLGlob("./dashboard/build/index.html")
		app.Static("/public", "dashboard/build")
	} else {
		logger.Info("Serve rice")
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
	dashboard.GET("/accessrole", context.HandleIndexPage(context_path))
	// login
	dashboard.GET("/login", context.HandleIndexPage(context_path))
	dashboard.POST("/login", context.HandleLogin(admin_user, admin_pass))
	// accessrole
	dashboard.GET("/role/account", context.HandleGetAllServiceAccount(&db))
	dashboard.POST("/role/account/:name", context.HandleCreateServiceAccount(&db))
}

func registerAPIRouter(app *gin.Engine, handler *httphandler.WebhookHandler, botClient *linebot.Client) {
	v1 := app.Group("/api/v1")
	// middleware
	v1.Use(middleware.UserAuthMiddleware(middleware.AllowPathPrefixSkipper(
		middleware.Prefix{Method: "Any", Path: "/api/v1/bot/hook"})))

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
