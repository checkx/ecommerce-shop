package entity

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestLoginReq_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		req     LoginReq
		wantErr bool
	}{
		{
			name: "valid login request",
			req: LoginReq{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			req: LoginReq{
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			req: LoginReq{
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			req: LoginReq{
				Email: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			req: LoginReq{
				Email:    "",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "empty password",
			req: LoginReq{
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegisterReq_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		req     RegisterReq
		wantErr bool
	}{
		{
			name: "valid register request",
			req: RegisterReq{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			req: RegisterReq{
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			req: RegisterReq{
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			req: RegisterReq{
				Email: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "password too short",
			req: RegisterReq{
				Email:    "test@example.com",
				Password: "123",
			},
			wantErr: true,
		},
		{
			name: "password exactly 8 characters",
			req: RegisterReq{
				Email:    "test@example.com",
				Password: "12345678",
			},
			wantErr: false,
		},
		{
			name: "empty email",
			req: RegisterReq{
				Email:    "",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "empty password",
			req: RegisterReq{
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthResponse_Structure(t *testing.T) {
	response := AuthResponse{
		ID:    "user-123",
		Token: "jwt-token-here",
	}

	assert.Equal(t, "user-123", response.ID)
	assert.Equal(t, "jwt-token-here", response.Token)
}
