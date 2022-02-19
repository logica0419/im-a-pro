package editor

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/nfnt/resize"
	"gocv.io/x/gocv"
)

func PutDetectedRect(rectangles []image.Rectangle, base image.Image) (image.Image, error) {
	mat, err := gocv.ImageToMatRGB(base)
	if err != nil {
		return nil, err
	}
	defer mat.Close()

	originRectangle := image.Rectangle{image.Point{0, 0}, base.Bounds().Size()}

	for _, rectangle := range rectangles {
		if rectangle.Max.X-rectangle.Min.X < originRectangle.Max.X/15 {
			continue
		}
		if rectangle.Max.Y-rectangle.Min.Y < originRectangle.Max.Y/15 {
			continue
		}

		gocv.Rectangle(&mat, rectangle, color.RGBA{0, 0, 255, 0}, 10)
	}

	newImg, err := mat.ToImage()
	if err != nil {
		return nil, err
	}

	return newImg, nil
}

func PutSun(rectangles []image.Rectangle, base image.Image) (image.Image, error) {
	sun, err := os.Open("./assets/sun.png")
	if err != nil {
		return nil, err
	}
	defer sun.Close()

	sunImg, _, err := image.Decode(sun)
	if err != nil {
		return nil, err
	}

	originRectangle := image.Rectangle{image.Point{0, 0}, base.Bounds().Size()}
	newImg := image.NewRGBA(originRectangle)
	draw.Draw(newImg, originRectangle, base, image.Point{0, 0}, draw.Src)

	for _, rectangle := range rectangles {
		if rectangle.Max.X-rectangle.Min.X < originRectangle.Max.X/15 {
			continue
		}
		if rectangle.Max.Y-rectangle.Min.Y < originRectangle.Max.Y/15 {
			continue
		}

		sunDrawArea := image.Rectangle{
			image.Point{
				rectangle.Min.X - (rectangle.Max.X-rectangle.Min.X)*4/10,
				rectangle.Min.Y - (rectangle.Max.Y-rectangle.Min.Y)*4/10,
			},
			image.Point{
				rectangle.Max.X + (rectangle.Max.X-rectangle.Min.X)*4/10,
				rectangle.Max.Y + (rectangle.Max.Y-rectangle.Min.Y)*4/10,
			},
		}
		resizedSun := resize.Resize(uint(sunDrawArea.Max.X-sunDrawArea.Min.X), uint(sunDrawArea.Max.Y-sunDrawArea.Min.Y), sunImg, resize.Lanczos3)

		draw.Draw(newImg, sunDrawArea, resizedSun, image.Point{0, 0}, draw.Over)
	}

	return newImg, nil
}
