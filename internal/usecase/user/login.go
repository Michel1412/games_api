package user

import (
	"errors"

	domainuser "games_api/internal/domain/user"

	"github.com/google/uuid"
)

const (
	validEmail    = "usuario@esoft.com"
	validPassword = "Abc123"
)

type LoginUseCase struct{}

func NewLoginUseCase() *LoginUseCase {
	return &LoginUseCase{}
}

func (uc *LoginUseCase) Execute(request domainuser.LoginRequest) (domainuser.LoginResponse, error) {
	if err := request.Validate(); err != nil {
		return domainuser.LoginResponse{}, err
	}

	if request.Email != validEmail || request.Password != validPassword {
		return domainuser.LoginResponse{}, errors.New("credenciais invalidas")
	}

	return domainuser.LoginResponse{Token: uuid.NewString()}, nil
}
