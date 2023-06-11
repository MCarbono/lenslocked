package services

import (
	"database/sql"
	"lenslocked/domain/entity"
	repository "lenslocked/repository/sqlite"
	"lenslocked/token"
	"os/exec"
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
	var passwordResetRepository = repository.NewPasswordResetSQLite(db)
	var userService = &UserService{
		UserRepository: userRepository,
	}
	var passwordResetService = &PasswordResetService{
		TokenManager:   token.ManagerImpl{},
		UserRepository: userRepository,
		PasswordReset:  passwordResetRepository,
	}
	type args struct {
		email            string
		passwordResetURL string
	}
	type test struct {
		name string
		args args
	}
	tests := []test{
		{
			name: "Should create a password reset and send an email to the user",
			args: args{
				email:            "teste@email.com",
				passwordResetURL: "http://localhost:3000/reset-pw?",
			},
		},
	}
	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			defer db.Exec("DELETE from users;")
			defer db.Exec("DELETE from password_resets;")
			_, err := userService.Create("teste@email.com", "password")
			if err != nil {
				t.Fatal(err)
			}
			got, err := passwordResetService.Create(scenario.args.email, scenario.args.passwordResetURL)
			if err != nil {
				t.Fatal(err)
			}
			want, err := passwordResetRepository.FindByID(got.ID)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(entity.PasswordReset{}, "Token")); diff != "" {
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
	_, err = db.Exec("DROP TABLE IF EXISTS password_resets;")
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
	_, err = db.Exec(`CREATE TABLE password_resets (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INT UNIQUE NOT NULL REFERENCES users (id) ON DELETE CASCADE, token_hash TEXT UNIQUE NOT NULL, expires_at TIMESTAMP NOT NULL);`)
	if err != nil {
		return nil, err
	}
	return db, nil
}