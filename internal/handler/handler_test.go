package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	domainjogo "games_api/internal/domain/jogo"
	"games_api/internal/service"
	jogousecase "games_api/internal/usecase/jogo"
	userusecase "games_api/internal/usecase/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRouter() http.Handler {
	repo := service.NewMemoryJogoRepository()

	return NewRouter(
		userusecase.NewLoginUseCase(),
		jogousecase.NewListJogosUseCase(repo),
		jogousecase.NewGetJogoUseCase(repo),
		jogousecase.NewCreateJogoUseCase(repo),
		jogousecase.NewUpdateJogoUseCase(repo),
		jogousecase.NewDeleteJogoUseCase(repo),
	)
}

func TestLoginHandlerReturnsToken(t *testing.T) {
	router := newTestRouter()
	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{
		"email":"usuario@esoft.com",
		"password":"Abc123"
	}`))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	var body map[string]string
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	assert.NotEmpty(t, body["token"])
}

func TestLoginHandlerReturnsBadRequestForInvalidCredentials(t *testing.T) {
	router := newTestRouter()
	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{
		"email":"usuario@esoft.com",
		"password":"errada"
	}`))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusBadRequest, response.Code)
	var body map[string]string
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	assert.Equal(t, "credenciais invalidas", body["error"])
}

func TestListJogosHandler(t *testing.T) {
	router := newTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/jogos", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	var body []domainjogo.JogoResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	require.Len(t, body, 2)
	assert.Equal(t, "The Legend of Zelda", body[0].Nome)
}

func TestCreateJogoHandlerReturnsCreated(t *testing.T) {
	router := newTestRouter()
	request := httptest.NewRequest(http.MethodPost, "/jogos", strings.NewReader(`{
		"nome":"Elden Ring",
		"tipo":"RPG",
		"nota":9,
		"review":"Desafiador e visualmente impecavel."
	}`))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusCreated, response.Code)
	var body domainjogo.JogoResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	assert.Equal(t, 3, body.ID)
	assert.Equal(t, "Elden Ring", body.Nome)
}

func TestCreateJogoHandlerReturnsBadRequestForInvalidNota(t *testing.T) {
	router := newTestRouter()
	request := httptest.NewRequest(http.MethodPost, "/jogos", strings.NewReader(`{
		"nome":"Elden Ring",
		"tipo":"RPG",
		"nota":11,
		"review":"Desafiador."
	}`))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusBadRequest, response.Code)
	var body map[string]string
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	assert.Equal(t, "nota deve estar entre 1 e 10", body["error"])
}

func TestDeleteJogoHandlerReturnsNoContent(t *testing.T) {
	router := newTestRouter()
	request := httptest.NewRequest(http.MethodDelete, "/jogos/1", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNoContent, response.Code)
	assert.Empty(t, response.Body.String())
}
