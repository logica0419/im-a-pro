package router

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/logica0419/im-a-pro/detector"
	"gocv.io/x/gocv"
)

func (r *Router) handleTextMessage(userID, replyToken string, mes *linebot.TextMessage) error {
	img, ok := imageCache[userID]
	if !ok {
		_, err := r.bot.ReplyMessage(replyToken, linebot.NewTextMessage("画像が必要です")).Do()
		if err != nil {
			return err
		}
		return fmt.Errorf("image not found")
	}
	defer delete(imageCache, userID)

	rectangles, err := detector.DetectFace(img)
	if err != nil {
		return err
	}

	switch mes.Text {
	case "Detect":
		mat, err := gocv.ImageToMatRGB(img)
		if err != nil {
			return err
		}
		defer mat.Close()

		for _, rectangle := range rectangles {
			gocv.Rectangle(&mat, rectangle, color.RGBA{0, 0, 255, 0}, 10)
		}

		newImg, err := mat.ToImage()
		if err != nil {
			return err
		}
		err = r.replyImage(newImg, userID, replyToken)
		if err != nil {
			return err
		}

	default:
		_, err := r.bot.ReplyMessage(replyToken, linebot.NewTextMessage(imageReply)).Do()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Router) replyImage(img image.Image, messageID, replyToken string) error {
	out, err := os.Create("./temp.jpg")
	if err != nil {
		fmt.Println(err)
	}
	defer os.Remove("./temp.jpg")
	defer out.Close()

	err = jpeg.Encode(out, img, nil)
	if err != nil {
		r.e.Logger.Print(err)
		return err
	}

	buffer, _ := os.Open("./temp.jpg")
	err = r.uploadImageToDrive(messageID+".jpg", buffer)
	if err != nil {
		return err
	}

	_, err = r.bot.ReplyMessage(replyToken, linebot.NewImageMessage(r.domain+"/"+messageID+".jpg", r.domain+"/"+messageID+".jpg")).Do()
	if err != nil {
		return err
	}

	return nil
}
