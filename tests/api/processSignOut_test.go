package api

import (
	"fmt"
	"io/ioutil"
	"lenslocked/infra/controllers"
	"lenslocked/infra/http/cookie"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/services"
	"lenslocked/tests/testinfra"
	"lenslocked/token"
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
	var userRepository = repository.NewUserRepositorySQLite(db)
	var sessionRepository = repository.NewSessionRepositorySQLite(db)
	var userService = &services.UserService{
		UserRepository: userRepository,
	}
	var sessionService = &services.SessionService{
		DB:                db,
		SessionRepository: sessionRepository,
		UserRepository:    userRepository,
		TokenManager:      token.ManagerImpl{},
	}
	var userController = controllers.Users{UserService: userService, SessionService: sessionService}
	_, err = userService.Create("teste@email.com", "password")
	if err != nil {
		t.Fatal(err)
	}
	r := testinfra.NewRouterTest(userController)
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
			_, err = sessionService.User(token)
			if diff := cmp.Diff(scenario.want, err.Error()); diff != "" {
				t.Errorf("ProcessSignOut mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
