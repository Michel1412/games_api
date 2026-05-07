package user

import (
	"testing"

	domainuser "games_api/internal/domain/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginUseCaseReturnsTokenForValidCredentials(t *testing.T) {
	uc := NewLoginUseCase()

	response, err := uc.Execute(domainuser.LoginRequest{
		Email:    "usuario@esoft.com",
		Password: "Abc123",
	})

	require.NoError(t, err)
	_, err = uuid.Parse(response.Token)
	assert.NoError(t, err)
}

func TestLoginUseCaseRejectsInvalidCredentials(t *testing.T) {
	uc := NewLoginUseCase()

	_, err := uc.Execute(domainuser.LoginRequest{
		Email:    "usuario@esoft.com",
		Password: "wrong",
	})

	require.Error(t, err)
	assert.Equal(t, "credenciais invalidas", err.Error())
}

func TestLoginUseCaseValidatesRequiredFields(t *testing.T) {
	uc := NewLoginUseCase()

	_, err := uc.Execute(domainuser.LoginRequest{})

	require.Error(t, err)
	assert.Equal(t, "email e obrigatorio", err.Error())
}
