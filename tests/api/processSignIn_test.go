package api

import (
	"fmt"
	"io/ioutil"
	"lenslocked/application/usecases"
	"lenslocked/domain/entity"
	"lenslocked/idGenerator"
	"lenslocked/infra/controllers"
	"lenslocked/infra/http/cookie"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/services"
	"lenslocked/tests/fakes"
	"lenslocked/tests/testinfra"
	"lenslocked/tokenManager"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	_ "github.com/mattn/go-sqlite3"
)

func TestProcessSignIn(t *testing.T) {
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
	var userRepository = repository.NewUserRepositorySQLite(db)
	var sessionRepository = repository.NewSessionRepositorySQLite(db)
	var userService = &services.UserService{
		UserRepository: userRepository,
		IDGenerator:    fakes.NewIDGeneratorFake(),
	}
	var sessionService = &services.SessionService{
		DB:                db,
		SessionRepository: sessionRepository,
		UserRepository:    userRepository,
		TokenManager:      tokenManager.New(),
		IDGenerator:       idGenerator.New(),
	}
	var userController = controllers.Users{UserService: userService, SessionService: sessionService}
	createUserUseCase := usecases.NewCreateUserUseCase(userRepository, fakes.NewIDGeneratorFake())
	_, err = createUserUseCase.Execute(&usecases.CreateUserInput{Email: "teste@email.com", Password: "password"})
	if err != nil {
		t.Fatal(err)
	}
	r := testinfra.NewRouterTest(userController, controllers.Galleries{})
	ts := httptest.NewServer(r)
	defer ts.Close()
	type args struct {
		email    string
		password string
	}
	type test struct {
		name string
		args args
		want *entity.User
	}
	tests := []test{
		{
			name: "Should process a sign in with a user that already registered",
			args: args{
				email:    "teste@email.com",
				password: "password",
			},
			want: &entity.User{
				ID:    "fakeUUID",
				Email: "teste@email.com",
			},
		},
	}
	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			data := url.Values{}
			data.Add("email", scenario.args.email)
			data.Add("password", scenario.args.password)
			jar, _ := cookiejar.New(nil)
			client := ts.Client()
			client.Jar = jar
			resp, err := client.Post(fmt.Sprintf("%s/signin", ts.URL), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				t.Errorf("ProcessSignIn request failed with error: %v", string(body))
				return
			}
			token, err := cookie.ReadCookie(resp.Request, cookie.CookieSession)
			if err != nil {
				t.Fatal(err)
			}
			user, err := sessionService.User(token)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(scenario.want, user, cmpopts.IgnoreFields(entity.User{}, "PasswordHash")); diff != "" {
				t.Errorf("ProcessSignIn mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
