package controllers

import (
	"database/sql"
	"fmt"
	"lenslocked/domain/entity"
	repository "lenslocked/repository/sqlite"
	"lenslocked/services"
	"lenslocked/token"
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

func TestCreate(t *testing.T) {
	t.Cleanup(func() {
		cmd := exec.Command("rm", "../lenslocked_test.db")
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	})
	db, err := createDatabaseTest()
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
	var userController = Users{UserService: userService, SessionService: sessionService}
	r := NewRouterTest(userController)
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
			name: "Should create a new user",
			args: args{
				email:    "teste@email.com",
				password: "password",
			},
			want: &entity.User{
				ID:    1,
				Email: "teste@email.com",
			},
		},
		{
			name: "Should create a new user with email in uppercase",
			args: args{
				email:    "TESTE@EMAIL.COM",
				password: "password",
			},
			want: &entity.User{
				ID:    2,
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
			resp, err := client.Post(fmt.Sprintf("%s/users", ts.URL), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			token, err := readCookie(resp.Request, CookieSession)
			if err != nil {
				t.Fatal(err)
			}
			user, err := sessionService.User(token)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(scenario.want, user, cmpopts.IgnoreFields(entity.User{}, "PasswordHash")); diff != "" {
				t.Errorf("Create mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestProcessSignIn(t *testing.T) {
	t.Cleanup(func() {
		cmd := exec.Command("rm", "../lenslocked_test.db")
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	})
	db, err := createDatabaseTest()
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
	var userController = Users{UserService: userService, SessionService: sessionService}
	_, err = userService.Create("teste@email.com", "password")
	if err != nil {
		t.Fatal(err)
	}
	r := NewRouterTest(userController)
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
			token, err := readCookie(resp.Request, CookieSession)
			if err != nil {
				t.Fatal(err)
			}
			user, err := sessionService.User(token)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(scenario.want, user, cmpopts.IgnoreFields(entity.User{}, "PasswordHash")); diff != "" {
				t.Errorf("Create mismatch (-want +got):\n%v", diff)
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
	db, err := createDatabaseTest()
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
	var userController = Users{UserService: userService, SessionService: sessionService}
	_, err = userService.Create("teste@email.com", "password")
	if err != nil {
		t.Fatal(err)
	}
	r := NewRouterTest(userController)
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
			token, err := readCookie(resp.Request, CookieSession)
			if err != nil {
				t.Fatal(err)
			}
			resp, err = client.Post(fmt.Sprintf("%s/signout", ts.URL), "application/x-www-form-urlencoded", nil)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			_, err = sessionService.User(token)
			if diff := cmp.Diff(scenario.want, err.Error()); diff != "" {
				t.Errorf("Create mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func createDatabaseTest() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "../lenslocked_test.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS users;")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS sessions;")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, email TEXT UNIQUE NOT NULL,password_hash TEXT NOT NULL);`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE sessions (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INT UNIQUE NOT NULL REFERENCES users (id) ON DELETE CASCADE, token_hash TEXT UNIQUE NOT NULL);`)
	if err != nil {
		return nil, err
	}
	return db, nil
}
