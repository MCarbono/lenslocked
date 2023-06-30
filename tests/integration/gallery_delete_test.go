package integration

import (
	"lenslocked/application/usecases"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/tests/fakes"
	"lenslocked/tests/testinfra"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDeleteGallery(t *testing.T) {
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
	var deleteGalleryUseCase = usecases.NewDeleteGalleryUseCase(galleryRepository)
	var findGalleryUseCase = usecases.NewFindGalleryUseCase(galleryRepository)

	type test struct {
		name  string
		input *usecases.CreateGalleryInput
		want  string
	}

	tests := []test{
		{
			name: "Should create a new gallery",
			input: &usecases.CreateGalleryInput{
				Title:  "Gallery fake test",
				UserID: "fakeUserID123",
			},
			want: "query gallery by ID sql: no rows in result set",
		},
	}
	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			got, err := createGalleryUseCase.Execute(scenario.input)
			if err != nil {
				t.Fatal(err)
			}
			err = deleteGalleryUseCase.Execute(got.ID)
			if err != nil {
				t.Fatal(err)
			}
			_, err = findGalleryUseCase.Execute(got.ID)
			if diff := cmp.Diff(scenario.want, err.Error()); diff != "" {
				t.Errorf("Delete Gallery mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
