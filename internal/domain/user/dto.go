package user

import (
	"errors"
	"strings"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	if strings.TrimSpace(r.Email) == "" {
		return errors.New("email e obrigatorio")
	}

	if strings.TrimSpace(r.Password) == "" {
		return errors.New("password e obrigatorio")
	}

	return nil
}

type LoginResponse struct {
	Token string `json:"token"`
}
