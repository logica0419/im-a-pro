package detector

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"gocv.io/x/gocv"
)

var classifier = gocv.NewCascadeClassifier()

func init() {
	if !classifier.Load(os.Getenv("CASCADE")) {
		fmt.Printf("Error reading cascade file")
		return
	}
}

func DetectFace(img image.Image) ([]image.Rectangle, error) {
	mat, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return nil, err
	}
	defer mat.Close()

	rectangles := classifier.DetectMultiScale(mat)

	return rectangles, nil
}
