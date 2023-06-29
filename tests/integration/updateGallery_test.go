package integration

import (
	"lenslocked/application/usecases"
	"lenslocked/domain/entity"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/tests/fakes"
	"lenslocked/tests/testinfra"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUpdateGallery(t *testing.T) {
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
	var galleryRepository = repository.NewGalleryRepositorySQLite(db)
	var createGalleryUseCase = usecases.NewCreateGalleryUseCase(galleryRepository, fakes.NewIDGeneratorFake())
	var updateGalleryUseCase = usecases.NewUpdateGalleryUseCase(galleryRepository)
	var findGalleryUseCase = usecases.NewFindGalleryUseCase(galleryRepository)

	type test struct {
		name               string
		createGalleryinput *usecases.CreateGalleryInput
		updateGalleryInput *usecases.UpdateGalleryInput
		want               *entity.Gallery
	}

	tests := []test{
		{
			name: "Should create a new gallery",
			createGalleryinput: &usecases.CreateGalleryInput{
				Title:  "Gallery fake test",
				UserID: "fakerUserID123",
			},
			updateGalleryInput: &usecases.UpdateGalleryInput{
				ID:    fakes.NewIDGeneratorFake().Generate(),
				Title: "Updated Gallery Title",
			},
			want: entity.NewGallery(fakes.NewIDGeneratorFake().Generate(), "fakerUserID123", "Updated Gallery Title"),
		},
	}
	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			galleryCreated, err := createGalleryUseCase.Execute(scenario.createGalleryinput)
			if err != nil {
				t.Fatal(err)
			}
			err = updateGalleryUseCase.Execute(scenario.updateGalleryInput)
			if err != nil {
				t.Fatal(err)
			}
			got, err := findGalleryUseCase.Execute(galleryCreated.ID)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(scenario.want, got); diff != "" {
				t.Errorf("Create Gallery mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
