package integration

import (
	"lenslocked/application/usecases"
	"lenslocked/domain/entity"
	repository "lenslocked/infra/repository/disk"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFindGalleries(t *testing.T) {
	type args struct {
		GalleryID string
		Filename  string
	}

	type test struct {
		name string
		args args
		want *entity.Image
	}

	var imageRepository = repository.NewImageRepositoryDisk("../assets/images/", []string{".png", ".jpg", ".jpeg", ".gif"})
	var findImageUseCase = usecases.NewFindImageUseCase(imageRepository)

	tests := []test{
		{
			name: "Should find an image",
			args: args{
				GalleryID: "fakeUUIDImages",
				Filename:  "IMG_9897.jpg",
			},
			want: &entity.Image{
				GalleryID:       "fakeUUIDImages",
				Path:            "../assets/images/gallery-fakeUUIDImages/IMG_9897.jpg",
				Filename:        "IMG_9897.jpg",
				FilenameEscaped: "IMG_9897.jpg",
			},
		},
		{
			name: "Should find an image with an '?' character and should escape it",
			args: args{
				GalleryID: "fakeUUIDImages",
				Filename:  "IMG_4?89.png",
			},
			want: &entity.Image{
				GalleryID:       "fakeUUIDImages",
				Path:            "../assets/images/gallery-fakeUUIDImages/IMG_4?89.png",
				Filename:        "IMG_4?89.png",
				FilenameEscaped: "IMG_4%3F89.png",
			},
		},
	}
	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			got, err := findImageUseCase.Execute(scenario.args.GalleryID, scenario.args.Filename)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(scenario.want, got); diff != "" {
				t.Errorf("Find Galleries mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
