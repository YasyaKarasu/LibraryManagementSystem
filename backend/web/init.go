package web

import (
	"fmt"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var e *echo.Echo

func initCors(e *echo.Echo) {
	corsConf := echo_middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderContentLength,
			echo.HeaderAccept,
		},
		AllowCredentials: true,
		MaxAge:           3600,
		ExposeHeaders: []string{
			echo.HeaderAuthorization,
			echo.HeaderSetCookie,
		},
	}
	e.Use(echo_middleware.CORSWithConfig(corsConf))
}

func InitWebFramework() {
	e = echo.New()
	e.HideBanner = true
	e.Use(echo_middleware.LoggerWithConfig(echo_middleware.LoggerConfig{
		Format: "method=${method}, remote_ip=${remote_ip}, uri=${uri}, status=${status}\n",
	}))

	initCors(e)
	addRoutes(e)

	logrus.Info("Echo framework initialized")
}

func StartServer() {
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", viper.GetInt("server.port"))))
}
