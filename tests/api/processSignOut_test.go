package api

import (
	"fmt"
	"io/ioutil"
	"lenslocked/application/usecases"
	"lenslocked/idGenerator"
	"lenslocked/infra/controllers"
	"lenslocked/infra/http/cookie"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/tests/assets/fakes"
	"lenslocked/tests/assets/testinfra"
	"lenslocked/tokenManager"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	_ "github.com/mattn/go-sqlite3"
)

func TestProcessSignOut(t *testing.T) {
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
	tokenManager := tokenManager.New()
	idGenerator := idGenerator.New()
	var userRepository = repository.NewUserRepositorySQLite(db)
	var sessionRepository = repository.NewSessionRepositorySQLite(db)
	var createSessionUseCase = usecases.NewCreateSessionUseCase(sessionRepository, tokenManager, idGenerator)
	var signInUseCase = usecases.NewSignInUseCase(sessionRepository, userRepository, tokenManager, idGenerator)
	var signOutUseCase = usecases.NewSignOutUseCase(sessionRepository, tokenManager)
	var findUserByTokenUseCase = usecases.NewFindUserByTokenUseCase(userRepository, tokenManager)
	var userController = controllers.Users{
		CreateSessionUseCase:   createSessionUseCase,
		SignInUseCase:          signInUseCase,
		SignOutUseCase:         signOutUseCase,
		FindUserByTokenUseCase: findUserByTokenUseCase,
	}
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
		want string
	}
	tests := []test{
		{
			name: "Should remove a session for a user that was signed in and did a log out",
			args: args{
				email:    "teste@email.com",
				password: "password",
			},
			want: "sql: no rows in result set",
		},
	}
	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			defer db.Exec("DELETE from users;")
			defer db.Exec("DELETE from sessions;")
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
				t.Errorf("ProcessSignOut request failed with error: %v", string(body))
				return
			}
			token, err := cookie.ReadCookie(resp.Request, cookie.CookieSession)
			if err != nil {
				t.Fatal(err)
			}
			resp, err = client.Post(fmt.Sprintf("%s/signout", ts.URL), "application/x-www-form-urlencoded", nil)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				t.Errorf("ProcessSignOut request failed with error: %v", string(body))
				return
			}
			_, err = findUserByTokenUseCase.Execute(token)
			if diff := cmp.Diff(scenario.want, err.Error()); diff != "" {
				t.Errorf("ProcessSignOut mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
