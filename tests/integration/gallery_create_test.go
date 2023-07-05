package integration

import (
	"lenslocked/application/usecases"
	repositoryDisk "lenslocked/infra/repository/disk"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/tests/assets/fakes"
	"lenslocked/tests/assets/testinfra"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateGallery(t *testing.T) {
	t.Cleanup(func() {
		cmd := exec.Command("rm", "../lenslocked_test.db")
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	})
	db, err := testinfra.CreateDatabaseTest()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var imageRepository = repositoryDisk.NewImageRepositoryDisk("../assets/images/", []string{".png", ".jpg", ".jpeg", ".gif"})
	var galleryRepository = repository.NewGalleryRepositorySQLite(db)
	var createGalleryUseCase = usecases.NewCreateGalleryUseCase(galleryRepository, fakes.NewIDGeneratorFake())
	var findGalleryUseCase = usecases.NewFindGalleryUseCase(galleryRepository, imageRepository)

	type test struct {
		name  string
		input *usecases.CreateGalleryInput
	}

	tests := []test{
		{
			name: "Should create a new gallery",
			input: &usecases.CreateGalleryInput{
				Title:  "Gallery fake test",
				UserID: "fakeUserID123",
			},
		},
	}
	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			got, err := createGalleryUseCase.Execute(scenario.input)
			if err != nil {
				t.Fatal(err)
			}
			want, err := findGalleryUseCase.Execute(got.ID)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("Create Gallery mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
