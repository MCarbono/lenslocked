package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"lenslocked/application/usecases"
	"lenslocked/infra/controllers"
	"lenslocked/infra/http/cookie"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/tests/assets/fakes"
	"lenslocked/tests/testinfra"
	"lenslocked/tokenManager"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
	_ "github.com/mattn/go-sqlite3"
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
	idGenerator := fakes.NewIDGeneratorFake()
	tokenManager := tokenManager.New()

	var userRepository = repository.NewUserRepositorySQLite(db)
	var galleryRepository = repository.NewGalleryRepositorySQLite(db)
	var sessionRepository = repository.NewSessionRepositorySQLite(db)
	var creteUserUseCase = usecases.NewCreateUserUseCase(userRepository, idGenerator)
	var findUserByTokenUseCase = usecases.NewFindUserByTokenUseCase(userRepository, tokenManager)
	var createGalleryUseCase = usecases.NewCreateGalleryUseCase(galleryRepository, idGenerator)
	var createSessionUseCase = usecases.NewCreateSessionUseCase(sessionRepository, tokenManager, idGenerator)
	var findGalleryUseCase = usecases.NewFindGalleryUseCase(galleryRepository)
	var deleteGalleryUseCase = usecases.NewDeleteGalleryUseCase(galleryRepository)
	var findGalleriesUseCase = usecases.NewFindGalleriesUseCase(galleryRepository)

	var userController = controllers.Users{
		FindUserByTokenUseCase: findUserByTokenUseCase,
	}

	var galleryController = controllers.Galleries{
		FindGalleryUseCase:   findGalleryUseCase,
		DeleteGalleryUseCase: deleteGalleryUseCase,
		FindGalleriesUseCase: findGalleriesUseCase,
	}

	r := testinfra.NewRouterTest(userController, galleryController)
	ts := httptest.NewServer(r)
	defer ts.Close()

	type test struct {
		name string
		want error
	}
	tests := []test{
		{
			name: "Should delete a new gallery",
			want: errors.New("query gallery by ID sql: no rows in result set"),
		},
	}
	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			user, err := creteUserUseCase.Execute(&usecases.CreateUserInput{Email: "teste@email.com", Password: "password"})
			if err != nil {
				t.Fatal(err)
			}
			session, err := createSessionUseCase.Execute(user.ID)
			if err != nil {
				t.Fatal(err)
			}
			gallery, err := createGalleryUseCase.Execute(&usecases.CreateGalleryInput{Title: "Gallery delete test", UserID: user.ID})
			if err != nil {
				t.Fatal(err)
			}
			cookie := cookie.NewCookie(cookie.CookieSession, session.Token)
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/galleries/%s/delete", ts.URL, gallery.ID), nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.AddCookie(cookie)
			resp, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				t.Errorf("Create request failed with error: %v", string(body))
				return
			}
			got, err := findGalleryUseCase.Execute(gallery.ID)
			if got != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(scenario.want.Error(), err.Error()); diff != "" {
				t.Errorf("Create gallery mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
