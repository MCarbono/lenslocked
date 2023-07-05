package integration

import (
	"lenslocked/application/usecases"
	"lenslocked/domain/entity"
	"lenslocked/gen/mock"
	"lenslocked/idGenerator"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/tests/assets/fakes"
	"lenslocked/tests/testinfra"
	"lenslocked/tokenManager"
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
	idGeneratorFake := fakes.NewIDGeneratorFake()
	createUserUseCase := usecases.NewCreateUserUseCase(userRepository, idGeneratorFake)

	type mockFields struct {
		mockEmailProvider *mock.MockEmailProvider
	}

	type args struct {
		input *usecases.ForgotPasswordInput
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
				input: &usecases.ForgotPasswordInput{
					Email:            "teste@email.com",
					ResetPasswordURL: "http://localhost:3000/reset-pw?",
				},
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
			forgotPasswordUseCase := usecases.NewForgotPasswordUseCase(
				userRepository,
				passwordResetRepository,
				f.mockEmailProvider,
				idGenerator.New(),
				tokenManager.New(),
			)
			_, err = createUserUseCase.Execute(&usecases.CreateUserInput{Email: "teste@email.com", Password: "password"})
			if err != nil {
				t.Fatal(err)
			}
			got, err := forgotPasswordUseCase.Execute(scenario.args.input)
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
