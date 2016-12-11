// Copyright © 2016 Jagger Wang <jaggerwang@gmail.com>

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"zqc/controllers"
	"zqc/middlewares"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server",
	Long:  `Run server.`,
	Run: func(cmd *cobra.Command, args []string) {
		addr := viper.GetString("server.listenAddr")
		uploadDir := viper.GetString("storage.local.dir")

		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.Mkdir(uploadDir, 0755)
		}

		e := echo.New()
		e.Debug = viper.GetBool("server.debug")
		e.HTTPErrorHandler = controllers.HttpErrorHandler
		e.ShutdownTimeout = 10 * time.Second

		initEchoLog(e)

		addMiddlewares(e)

		addRoutes(e)

		e.Logger.Info("server pid ", os.Getpid())
		e.Logger.Info("server listening on ", addr)
		e.Logger.Fatal(e.Start(addr))
	},
}

func init() {
	serverCmd.Flags().StringP("server.listenAddr", "l", "", "server listen address")

	viper.BindPFlags(serverCmd.PersistentFlags())

	viper.BindPFlags(serverCmd.Flags())

	valid.SetFieldsRequiredByDefault(true)
}

func initEchoLog(e *echo.Echo) {
	w, err := os.OpenFile(filepath.Join(viper.GetString("dir.data"), viper.GetString("log.echo.file")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		panic(err)
	}
	e.Logger.SetOutput(w)

	var lvl log.Lvl
	switch strings.ToUpper(viper.GetString("log.level")) {
	case "DEBUG":
		lvl = log.DEBUG
	case "INFO":
		lvl = log.INFO
	case "WARN":
		lvl = log.WARN
	case "ERROR":
		lvl = log.ERROR
	case "OFF":
		lvl = log.OFF
	default:
		lvl = log.INFO
	}
	e.Logger.SetLevel(lvl)
}

func addMiddlewares(e *echo.Echo) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.Recover())

	e.Use(middleware.BodyLimit("100MB"))

	lcfg := middleware.DefaultLoggerConfig
	w, err := os.OpenFile(filepath.Join(viper.GetString("dir.data"), viper.GetString("log.request.file")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		panic(err)
	}
	lcfg.Output = w
	e.Use(middleware.LoggerWithConfig(lcfg))

	e.Use(middlewares.Session())

	e.Use(middlewares.MiddlewareContext())
}

func addRoutes(e *echo.Echo) {
	auth := middlewares.Auth()

	e.Static("/upload", viper.GetString("storage.local.dir"))

	e.POST("/register", controllers.RegisterAccount)
	e.GET("/login", controllers.Login)
	e.GET("/resetPassword", controllers.ResetPassword)
	e.GET("/isLogined", controllers.IsLogined)
	e.GET("/logout", controllers.Logout)

	var g *echo.Group

	g = e.Group("/security")
	g.GET("/sendVerifyCode", controllers.SendVerifyCode)

	g = e.Group("/account", auth)
	g.POST("/edit", controllers.EditAccount)
	g.GET("/info", controllers.AccountInfo)

	g = e.Group("/user", auth)
	g.GET("/info", controllers.UserInfo)
	g.GET("/nearby", controllers.NearbyUsers)

	g = e.Group("/file", auth)
	g.POST("/upload", controllers.UploadFile)
	g.GET("/info", controllers.FileInfo)
}
