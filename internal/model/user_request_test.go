package model

import "testing"

func TestCreateUserRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		request CreateUserRequest
		want    string
	}{
		{
			name: "valid request",
			request: CreateUserRequest{
				Username: "eko",
				Password: "rahasia123",
				Role:     "viewer",
				Email:    "eko@example.com",
			},
			want: "",
		},
		{
			name: "username required",
			request: CreateUserRequest{
				Username: "",
				Password: "rahasia123",
				Role:     "viewer",
				Email:    "eko@example.com",
			},
			want: "Username is required",
		},
		{
			name: "password minimum length",
			request: CreateUserRequest{
				Username: "eko",
				Password: "123",
				Role:     "viewer",
				Email:    "eko@example.com",
			},
			want: "Password must be at least 8 characters",
		},
		{
			name: "invalid role",
			request: CreateUserRequest{
				Username: "eko",
				Password: "rahasia123",
				Role:     "superadmin",
				Email:    "eko@example.com",
			},
			want: "Invalid role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.Validate()

			if got != tt.want {
				t.Fatalf(
					"expected %q, got %q",
					tt.want,
					got,
				)
			}
		})
	}
}
