package handler

import (
	"encoding/json"
	"net/http"

	domainuser "games_api/internal/domain/user"
)

type UserHandler struct {
	loginUC LoginUseCase
}

func NewUserHandler(loginUC LoginUseCase) *UserHandler {
	return &UserHandler{loginUC: loginUC}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request domainuser.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeBadRequest(w, "json invalido")
		return
	}

	response, err := h.loginUC.Execute(request)
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, response)
}
