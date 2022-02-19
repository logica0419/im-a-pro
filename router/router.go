package router

import (
	"context"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"google.golang.org/api/drive/v3"
)

type Router struct {
	e      *echo.Echo
	bot    *linebot.Client
	drive  *drive.Service
	token  string
	domain string
}

func NewRouter() (*Router, error) {
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("ACCESS_TOKEN"))
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	drive, err := drive.NewService(ctx)
	if err != nil {
		log.Panic(err)
	}

	e := newEcho()

	r := &Router{
		e:      e,
		bot:    bot,
		drive:  drive,
		token:  os.Getenv("ACCESS_TOKEN"),
		domain: os.Getenv("DOMAIN"),
	}

	r.e.POST("/", r.handleLineEvent)
	r.e.GET("/:imageName", r.handleImage)

	err = r.initIDCache()
	if err != nil {
		return nil, err
	}

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
