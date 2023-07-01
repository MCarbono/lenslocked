package api

import (
	"fmt"
	"io/ioutil"
	"lenslocked/application/usecases"
	"lenslocked/domain/entity"
	"lenslocked/idGenerator"
	"lenslocked/infra/controllers"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/rand"
	"lenslocked/services"
	"lenslocked/tests/testinfra"
	"lenslocked/tokenManager"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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
	idGenerator := idGenerator.New()
	var userRepository = repository.NewUserRepositorySQLite(db)
	var sessionRepository = repository.NewSessionRepositorySQLite(db)
	var passwordResetRepository = repository.NewPasswordResetSQLite(db)
	var userService = &services.UserService{
		UserRepository: userRepository,
		IDGenerator:    idGenerator,
	}
	createUserUseCase := usecases.NewCreateUserUseCase(userRepository, idGenerator)
	var sessionService = &services.SessionService{
		DB:                db,
		SessionRepository: sessionRepository,
		UserRepository:    userRepository,
		TokenManager:      tokenManager.New(),
		IDGenerator:       idGenerator,
	}

	var passwordResetService = &services.PasswordResetService{
		TokenManager:      tokenManager.New(),
		PasswordReset:     passwordResetRepository,
		UserRepository:    userRepository,
		SessionRepository: sessionRepository,
		IDGenerator:       idGenerator,
	}
	var userController = controllers.Users{PasswordResetService: passwordResetService, SessionService: sessionService, UserService: userService}
	r := testinfra.NewRouterTest(userController, controllers.Galleries{})
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
			user, err := createUserUseCase.Execute(&usecases.CreateUserInput{Email: "teste@email.com", Password: "password"})
			if err != nil {
				t.Fatal(err)
			}
			token, err := rand.String(32)
			if err != nil {
				t.Fatal(err)
			}
			tokenHash := passwordResetService.TokenManager.Hash(token)
			passwordReset := entity.NewPasswordReset(passwordResetService.Generate(), user.ID, token, tokenHash, 1*time.Hour)
			err = passwordResetService.PasswordReset.Create(passwordReset)
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
			passwordResetDeleted, err := passwordResetRepository.FindByID(passwordReset.ID)
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
