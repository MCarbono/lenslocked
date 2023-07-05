package usecases

import (
	"fmt"
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
	"path/filepath"
	"strings"
)

type FindGalleryUseCase struct {
	galleryRepository repository.GalleryRepository
	imagesDir         string
}

func NewFindGalleryUseCase(galleryRepository repository.GalleryRepository) *FindGalleryUseCase {
	return &FindGalleryUseCase{
		galleryRepository: galleryRepository,
	}
}

func (uc *FindGalleryUseCase) Execute(ID string) (*entity.Gallery, error) {
	gallery, err := uc.galleryRepository.FindByID(ID)
	if err != nil {
		return nil, fmt.Errorf("query gallery by ID %w", err)
	}
	globPattern := filepath.Join(uc.galleryDir(ID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}
	var images []*entity.Image
	for _, file := range allFiles {
		if uc.hasExtension(file, uc.extensions()) {
			images = append(images, &entity.Image{
				Path:      file,
				GalleryID: gallery.ID,
				Filename:  filepath.Base(file),
			})
		}
	}
	gallery.AddImages(images)
	return gallery, nil
}

func (uc *FindGalleryUseCase) galleryDir(ID string) string {
	imagesDir := uc.imagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%s", ID))
}

func (uc *FindGalleryUseCase) hasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}

func (uc *FindGalleryUseCase) extensions() []string {
	return []string{".png", ".jpg", ".jpeg", ".gif"}
}
