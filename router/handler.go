package router

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (r *Router) handleLineEvent(c echo.Context) error {
	events, err := r.bot.ParseRequest(c.Request())
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			r.e.Logger.Print(err)
		}
		return err
	}

	log.Printf("%#v", events)

	return nil
}
