package router

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Router struct {
	e   *echo.Echo
	bot *linebot.Client
}

func NewRouter() (*Router, error) {
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("ACCESS_TOKEN"))
	if err != nil {
		return nil, err
	}

	e := newEcho()
	r := &Router{
		e:   e,
		bot: bot,
	}

	r.e.POST("/", r.handleLineEvent)

	return r, nil
}

func newEcho() *echo.Echo {
	e := echo.New()

	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetHeader("${time_rfc3339} ${prefix} ${short_file} ${line} |")
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: "${time_rfc3339} method = ${method} | uri = ${uri} | status = ${status} ${error}\n"}))

	return e
}

func (r *Router) Run() {
	r.e.Logger.Panic(r.e.Start(":8000"))
}
