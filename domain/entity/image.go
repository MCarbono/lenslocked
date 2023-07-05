package entity

import (
	"net/url"
	"path/filepath"
)

type Image struct {
	GalleryID       string
	Path            string
	Filename        string
	FilenameEscaped string
}

func NewImage(galleryID, path, filename string) *Image {
	return &Image{
		GalleryID:       galleryID,
		Path:            path,
		Filename:        filename,
		FilenameEscaped: url.PathEscape(filepath.Base(filename)),
	}
}
