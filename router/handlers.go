package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (r *Router) handleImage(e echo.Context) error {
	imageName := e.Param("imageName")

	img, err := r.getImageFromDrive(imageName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return e.Stream(http.StatusOK, "image/jpeg", img)
}

func (r *Router) handleLineEvent(c echo.Context) error {
	events, err := r.bot.ParseRequest(c.Request())
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			r.e.Logger.Print(err)
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				err := r.handleTextMessage(event.Source.UserID, event.ReplyToken, message)
				if err != nil {
					r.e.Logger.Print(err)
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
			case *linebot.ImageMessage:
				err = r.handleImageMessage(event.Source.UserID, event.ReplyToken, message)
				if err != nil {
					r.e.Logger.Print(err)
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
			}
		}
	}

	return c.NoContent(http.StatusOK)
}
