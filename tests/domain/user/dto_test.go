package user_test

import (
	"testing"

	domainuser "games_api/internal/domain/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginRequestValidate(t *testing.T) {
	cases := []struct {
		name    string
		request domainuser.LoginRequest
		wantErr string
	}{
		{
			name:    "valido",
			request: domainuser.LoginRequest{Email: "user@esoft.com", Password: "Abc123"},
		},
		{
			name:    "email vazio",
			request: domainuser.LoginRequest{Email: "  ", Password: "Abc123"},
			wantErr: "email e obrigatorio",
		},
		{
			name:    "password vazia",
			request: domainuser.LoginRequest{Email: "user@esoft.com", Password: " "},
			wantErr: "password e obrigatorio",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.request.Validate()
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}

			require.Error(t, err)
			assert.Equal(t, tc.wantErr, err.Error())
		})
	}
}
