package integration

import (
	"lenslocked/domain/entity"
	"lenslocked/gen/mock"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/services"
	"lenslocked/tests/testinfra"
	"lenslocked/token"
	"os/exec"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	_ "github.com/mattn/go-sqlite3"
)

func TestCreatePasswordReset(t *testing.T) {
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
	var passwordResetRepository = repository.NewPasswordResetSQLite(db)
	var userService = &services.UserService{
		UserRepository: userRepository,
	}
	var passwordResetService = &services.PasswordResetService{
		TokenManager:   token.ManagerImpl{},
		UserRepository: userRepository,
		PasswordReset:  passwordResetRepository,
	}

	type mockFields struct {
		mockEmailProvider *mock.MockEmailProvider
	}

	type args struct {
		email            string
		passwordResetURL string
	}

	type test struct {
		name  string
		args  args
		mocks func(f *mockFields)
	}
	tests := []test{
		{
			name: "Should create a password reset and send an email to the user",
			args: args{
				email:            "teste@email.com",
				passwordResetURL: "http://localhost:3000/reset-pw?",
			},
			mocks: func(f *mockFields) {
				f.mockEmailProvider.EXPECT().Send(gomock.Any()).Return(nil).Times(1)
			},
		},
	}
	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			defer db.Exec("DELETE from users;")
			defer db.Exec("DELETE from password_resets;")
			mockCtrl := gomock.NewController(t)
			f := mockFields{
				mockEmailProvider: mock.NewMockEmailProvider(mockCtrl),
			}
			scenario.mocks(&f)
			passwordResetService.EmailGateway = f.mockEmailProvider
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
