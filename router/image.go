package router

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var imageCache = map[string]image.Image{}

const imageReply = "次のうち、どの処理を行いますか？\n" + "・Detect (顔の検出)\n" + "・Sun (顔に太陽を付ける)"

func (r *Router) handleImageMessage(userID, replyToken string, mes *linebot.ImageMessage) error {
	if _, ok := imageCache[mes.ID]; ok {
		_, err := r.bot.ReplyMessage(replyToken, linebot.NewTextMessage("既に画像が登録されています")).Do()
		if err != nil {
			return err
		}

		return fmt.Errorf("image already exists")
	}

	img, err := r.getImageFromLine(mes.ID)
	if err != nil {
		return err
	}
	defer img.Close()

	decoded, _, err := image.Decode(img)
	if err != nil {
		return err
	}

	imageCache[userID] = decoded

	_, err = r.bot.ReplyMessage(replyToken, linebot.NewTextMessage(imageReply)).Do()
	if err != nil {
		return err
	}

	return nil
}

func (r *Router) getImageFromLine(messageID string) (io.ReadCloser, error) {
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
