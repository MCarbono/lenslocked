package integration

import (
	"lenslocked/application/usecases"
	"lenslocked/domain/entity"
	"lenslocked/idGenerator"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/tests/assets/testinfra"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFindGalleries(t *testing.T) {
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
	var createGalleryUseCase = usecases.NewCreateGalleryUseCase(galleryRepository, idGenerator.New())
	var findGalleriesUseCase = usecases.NewFindGalleriesUseCase(galleryRepository)

	type test struct {
		name  string
		input []*usecases.CreateGalleryInput
		want  []*entity.Gallery
	}

	tests := []test{
		{
			name: "Should create two galleries and return find it",
			input: []*usecases.CreateGalleryInput{
				{
					Title:  "Gallery fake test one",
					UserID: "fakeUserID123",
				},
				{
					Title:  "Gallery fake test two",
					UserID: "fakeUserID123",
				},
			},
			want: []*entity.Gallery{
				{
					UserID: "fakeUserID123",
					Title:  "Gallery fake test one",
				}, {
					UserID: "fakeUserID123",
					Title:  "Gallery fake test two",
				},
			},
		},
	}
	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			for _, v := range scenario.input {
				_, err := createGalleryUseCase.Execute(v)
				if err != nil {
					t.Fatal(err)
				}
			}
			got, err := findGalleriesUseCase.Execute(scenario.input[0].UserID)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(scenario.want, got, cmpopts.IgnoreFields(entity.Gallery{}, "ID")); diff != "" {
				t.Errorf("Find Galleries mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
