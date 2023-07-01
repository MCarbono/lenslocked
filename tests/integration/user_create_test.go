package integration

import (
	"lenslocked/application/usecases"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/tests/fakes"
	"lenslocked/tests/testinfra"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateUser(t *testing.T) {
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
	var creteUserUseCase = usecases.NewCreateUserUseCase(userRepository, fakes.NewIDGeneratorFake())

	type test struct {
		name  string
		input *usecases.CreateUserInput
	}

	tests := []test{
		{
			name: "Should create a new user",
			input: &usecases.CreateUserInput{
				Email:    "user@email.com",
				Password: "password",
			},
		},
		{
			name: "Should create a new user with email as uppercase",
			input: &usecases.CreateUserInput{
				Email:    "USER@EMAIL.COM",
				Password: "password",
			},
		},
	}
	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			defer db.Exec("DELETE from users;")
			got, err := creteUserUseCase.Execute(scenario.input)
			if err != nil {
				t.Fatal(err)
			}
			want, err := userRepository.FindByID(got.ID)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("Create Gallery mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
