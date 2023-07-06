package integration

import (
	"io/fs"
	"lenslocked/application/usecases"
	repository "lenslocked/infra/repository/disk"
	"os/exec"
	"syscall"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDeleteImage(t *testing.T) {

	type args struct {
		GalleryID string
		Filename  string
	}

	type test struct {
		name string
		args args
		want *fs.PathError
	}

	var imageRepository = repository.NewImageRepositoryDisk("../assets/images/", []string{".png", ".jpg", ".jpeg", ".gif"})
	var deleteImageUseCase = usecases.NewDeleteImageUseCase(imageRepository)
	var findImageUseCase = usecases.NewFindImageUseCase(imageRepository)

	tests := []test{
		{
			name: "Should delete an image",
			args: args{
				GalleryID: "fakeUUIDImages",
				Filename:  "delete-test-image.png",
			},
			want: &fs.PathError{
				Op:   "stat",
				Path: "../assets/images/gallery-fakeUUIDImages/delete-test-image.png",
				Err:  syscall.Errno(0x02),
			},
		},
	}

	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			cmd := exec.Command("cp", "../assets/images/gallery-fakeUUIDImages/IMG_4?89.png", "../assets/images/gallery-fakeUUIDImages/delete-test-image.png")
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			err = deleteImageUseCase.Execute(scenario.args.GalleryID, scenario.args.Filename)
			if err != nil {
				t.Fatal(err)
			}
			image, err := findImageUseCase.Execute(scenario.args.GalleryID, scenario.args.Filename)
			if image != nil {
				t.Fatalf("the findImageUseCase() should not return a image, but got %v", image)
			}
			if diff := cmp.Diff(scenario.want, err); diff != "" {
				t.Errorf("Delete Image mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
