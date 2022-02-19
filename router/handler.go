package router

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

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
				_, err = r.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do()
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

func (r *Router) handleImageMessage(userID, replyToken string, mes *linebot.ImageMessage) error {
	img, err := r.getImage(mes.ID)
	if err != nil {
		return err
	}
	defer img.Close()

	file, err := os.Create("./images/" + mes.ID + ".jpg")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, img)
	if err != nil {
		return err
	}

	_, err = r.bot.ReplyMessage(replyToken, linebot.NewImageMessage(r.domain+"/"+mes.ID+".jpg", r.domain+"/"+mes.ID+".jpg")).Do()
	if err != nil {
		return err
	}

	return nil
}

func (r *Router) getImage(messageID string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", "https://api-data.line.me/v2/bot/message/"+messageID+"/content", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+r.token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
