package api

import (
	"fmt"
	"io/ioutil"
	"lenslocked/application/usecases"
	"lenslocked/domain/entity"
	"lenslocked/infra/controllers"
	"lenslocked/infra/http/cookie"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/tests/fakes"
	"lenslocked/tests/testinfra"
	"lenslocked/tokenManager"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	_ "github.com/mattn/go-sqlite3"
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
	idGenerator := fakes.NewIDGeneratorFake()
	tokenManager := tokenManager.New()

	var userRepository = repository.NewUserRepositorySQLite(db)
	var galleryRepository = repository.NewGalleryRepositorySQLite(db)
	var sessionRepository = repository.NewSessionRepositorySQLite(db)
	var creteUserUseCase = usecases.NewCreateUserUseCase(userRepository, idGenerator)
	var findUserByTokenUseCase = usecases.NewFindUserByTokenUseCase(userRepository, tokenManager)
	var createGalleryUseCase = usecases.NewCreateGalleryUseCase(galleryRepository, idGenerator)
	var createSessionUseCase = usecases.NewCreateSessionUseCase(sessionRepository, tokenManager, idGenerator)
	var findGalleriesUseCase = usecases.NewFindGalleryUseCase(galleryRepository)
	var updateGalleryUseCase = usecases.NewUpdateGalleryUseCase(galleryRepository)

	var userController = controllers.Users{
		CreateUserUseCase:      creteUserUseCase,
		FindUserByTokenUseCase: findUserByTokenUseCase,
	}

	var galleryController = controllers.Galleries{
		CreateGalleryUseCase: createGalleryUseCase,
		FindGalleryUseCase:   findGalleriesUseCase,
		UpdateGalleryUseCase: updateGalleryUseCase,
	}

	r := testinfra.NewRouterTest(userController, galleryController)
	ts := httptest.NewServer(r)
	defer ts.Close()

	type args struct {
		title string
	}
	type test struct {
		name string
		args args
		want *entity.Gallery
	}
	tests := []test{
		{
			name: "Should update a new gallery",
			args: args{
				title: "Gallery title test updated",
			},
			want: entity.NewGallery(idGenerator.Generate(), idGenerator.Generate(), "Gallery title test updated"),
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
			gallery, err := createGalleryUseCase.Execute(&usecases.CreateGalleryInput{Title: "Galley test", UserID: user.ID})
			if err != nil {
				t.Fatal(err)
			}
			cookie := cookie.NewCookie(cookie.CookieSession, session.Token)
			data := url.Values{}
			data.Add("title", scenario.args.title)
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/galleries/%s", ts.URL, gallery.ID), strings.NewReader(data.Encode()))
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
			url, err := resp.Request.Response.Location()
			if err != nil {
				t.Fatal(err)
			}
			galleryID := strings.Split(url.Path, "/")[2]
			got, err := findGalleriesUseCase.Execute(galleryID)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(scenario.want, got); diff != "" {
				t.Errorf("Update gallery mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
