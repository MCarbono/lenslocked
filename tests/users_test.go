package tests

import (
	"fmt"
	"io/ioutil"
	"lenslocked/domain/entity"
	"lenslocked/infra/controllers"
	"lenslocked/infra/http/cookie"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/rand"
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
	"time"

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
				ID:    1,
				Email: "teste@email.com",
			},
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

func TestProcessResetPassword(t *testing.T) {
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
	var passwordResetRepository = repository.NewPasswordResetSQLite(db)
	var userService = &services.UserService{
		UserRepository: userRepository,
	}
	var sessionService = &services.SessionService{
		DB:                db,
		SessionRepository: sessionRepository,
		UserRepository:    userRepository,
		TokenManager:      token.ManagerImpl{},
	}

	var passwordResetService = &services.PasswordResetService{
		TokenManager:      token.ManagerImpl{},
		PasswordReset:     passwordResetRepository,
		UserRepository:    userRepository,
		SessionRepository: sessionRepository,
	}
	var userController = controllers.Users{PasswordResetService: passwordResetService, SessionService: sessionService, UserService: userService}
	r := testinfra.NewRouterTest(userController)
	ts := httptest.NewServer(r)
	defer ts.Close()

	type args struct {
		newPassword string
	}

	type test struct {
		name string
		args args
	}

	tests := []test{
		{
			name: "Should update the user password.",
			args: args{
				newPassword: "123456",
			},
		},
	}
	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			defer db.Exec("DELETE from users;")
			defer db.Exec("DELETE from sessions;")
			defer db.Exec("DELETE from password_resets;")
			user, err := userService.Create("teste@email.com", "password")
			if err != nil {
				t.Fatal(err)
			}
			token, err := rand.String(32)
			if err != nil {
				t.Fatal(err)
			}
			tokenHash := passwordResetService.TokenManager.Hash(token)
			passwordReset := entity.NewPasswordReset(user.ID, token, tokenHash, 1*time.Hour)
			passwordResetID, err := passwordResetService.PasswordReset.Create(passwordReset)
			if err != nil {
				t.Fatal(err)
			}
			data := url.Values{}
			data.Add("password", scenario.args.newPassword)
			data.Add("token", passwordReset.Token)
			jar, _ := cookiejar.New(nil)
			client := ts.Client()
			client.Jar = jar
			resp, err := client.Post(fmt.Sprintf("%s/reset-pw", ts.URL), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				t.Errorf("Request failed with error: %v", string(body))
				return
			}
			userChangedPassword, err := userRepository.FindByID(user.ID)
			if err != nil {
				t.Fatal(err)
			}
			if userChangedPassword.PasswordHash == user.PasswordHash {
				t.Errorf("ProcessForgotPassword failed. User password should be updated")
				return
			}
			passwordResetDeleted, err := passwordResetRepository.FindByID(passwordResetID)
			if passwordResetDeleted != nil && err == nil {
				t.Errorf("ProcessForgotPassword failed. Entity PasswordReset should be deleted, but it's stored in the database. %v", passwordResetDeleted)
				return
			}
			wantErr := "password_resets: sql: no rows in result set"
			if err.Error() != wantErr {
				t.Errorf("ProcessForgotPassword failed. Err want %v\n got %v", wantErr, err)
			}
		})
	}
}
