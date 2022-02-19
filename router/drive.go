package router

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/api/drive/v3"
)

func (r *Router) handleImage(e echo.Context) error {
	imageName := e.Param("imageName")
	imageID := imageIDCache[imageName]

	if imageID == "" {
		return echo.NewHTTPError(http.StatusNotFound, "Image not found")
	}

	res, err := r.drive.Files.Get(imageID).Download()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return e.Stream(http.StatusOK, "image/jpeg", res.Body)
}

var imageIDCache map[string]string

func (r *Router) initIDCache() error {
	imageIDCache = make(map[string]string)

	fileList, err := r.drive.Files.List().Fields("files(id, name, mimeType)").Do()
	if err != nil {
		return err
	}

	for _, file := range fileList.Files {
		if file.MimeType == "application/vnd.google-apps.folder" {
			continue
		}

		imageIDCache[file.Name] = file.Id
	}

	return nil
}

func (r *Router) uploadImageToDrive(name string, image io.Reader) error {
	res, err := r.drive.Files.Create(&drive.File{Name: name}).Media(image).Do()
	if err != nil {
		return err
	}

	imageIDCache[name] = res.Id
	return nil
}
