package router

import (
	"fmt"
	"io"

	"google.golang.org/api/drive/v3"
)

var imageIDCache = map[string]string{}

func (r *Router) initIDCache() error {
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

func (r *Router) getImageFromDrive(imageName string) (io.Reader, error) {
	imageID := imageIDCache[imageName]

	if imageID == "" {
		return nil, fmt.Errorf("image not found")
	}

	res, err := r.drive.Files.Get(imageID).Download()
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func (r *Router) uploadImageToDrive(name string, image io.Reader) error {
	res, err := r.drive.Files.Create(&drive.File{Name: name}).Media(image).Do()
	if err != nil {
		return err
	}

	imageIDCache[name] = res.Id
	return nil
}
