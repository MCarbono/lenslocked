package entity

import (
	"lenslocked/domain/entity"
	"testing"
	"time"
)

func TestPasswordReset_IsExpired(t *testing.T) {
	type fields struct {
		ID        string
		UserID    string
		Token     string
		TokenHash string
		Duration  time.Duration
	}

	type test struct {
		name   string
		fields fields
		want   bool
	}

	tests := []test{
		{
			name: "The password reset should not be expired",
			fields: fields{
				ID:        "fakeID",
				UserID:    "fakeUserID",
				Token:     "fakeToken",
				TokenHash: "tokenHashFake123",
				Duration:  1 * time.Hour,
			},
			want: false,
		},
		{
			name: "The password reset should  be expired",
			fields: fields{
				ID:        "fakeID",
				UserID:    "fakeUserID",
				Token:     "fakeToken",
				TokenHash: "tokenHashFake123",
				Duration:  -1 * time.Hour,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pw := entity.NewPasswordReset(tt.fields.ID, tt.fields.UserID, tt.fields.Token, tt.fields.TokenHash, tt.fields.Duration)
			if got := pw.IsExpired(); got != tt.want {
				t.Errorf("PasswordReset.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}
