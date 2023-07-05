package integration

import (
	"lenslocked/application/usecases"
	"lenslocked/domain/entity"
	repository "lenslocked/infra/repository/sqlite"
	"lenslocked/tests/assets/fakes"
	"lenslocked/tests/testinfra"
	"lenslocked/tokenManager"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreateSession(t *testing.T) {
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
	var sessionRepository = repository.NewSessionRepositorySQLite(db)
	var createSessionUseCase = usecases.NewCreateSessionUseCase(sessionRepository, tokenManager.New(), fakes.NewIDGeneratorFake())

	type test struct {
		name   string
		userID string
	}

	tests := []test{
		{
			name:   "Should create a new session",
			userID: "fakeUserID123",
		},
	}
	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			got, err := createSessionUseCase.Execute(scenario.userID)
			if err != nil {
				t.Fatal(err)
			}
			want, err := sessionRepository.FindByTokenHash(got.TokenHash)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(entity.Session{}, "Token")); diff != "" {
				t.Errorf("Create Gallery mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
