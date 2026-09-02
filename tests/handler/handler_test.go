package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	domainjogo "games_api/internal/domain/jogo"
	domainwebhook "games_api/internal/domain/webhook"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestCreateWebhookHandlerReturnsCreated(t *testing.T) {
	router := newTestRouter()
	request := httptest.NewRequest(http.MethodPost, "/webhooks", strings.NewReader(`{
		"url":"https://example.com/hook"
	}`))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusCreated, response.Code)
	var body domainwebhook.WebhookResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	assert.Equal(t, 1, body.ID)
	assert.Equal(t, "https://example.com/hook", body.URL)
	assert.True(t, body.Ativo)
}

func TestCreateWebhookHandlerReturnsBadRequestForInvalidURL(t *testing.T) {
	router := newTestRouter()
	request := httptest.NewRequest(http.MethodPost, "/webhooks", strings.NewReader(`{"url":"ftp://example.com"}`))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusBadRequest, response.Code)
	var body map[string]string
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	assert.Equal(t, "url deve usar http ou https", body["error"])
}

func TestDeactivateAndActivateWebhookHandler(t *testing.T) {
	router := newTestRouter()

	createReq := httptest.NewRequest(http.MethodPost, "/webhooks", strings.NewReader(`{"url":"https://example.com/hook"}`))
	createRes := httptest.NewRecorder()
	router.ServeHTTP(createRes, createReq)
	require.Equal(t, http.StatusCreated, createRes.Code)

	deactivateReq := httptest.NewRequest(http.MethodPost, "/webhooks/1/desativar", nil)
	deactivateRes := httptest.NewRecorder()
	router.ServeHTTP(deactivateRes, deactivateReq)

	require.Equal(t, http.StatusOK, deactivateRes.Code)
	var deactivated domainwebhook.WebhookResponse
	require.NoError(t, json.NewDecoder(deactivateRes.Body).Decode(&deactivated))
	assert.False(t, deactivated.Ativo)

	activateReq := httptest.NewRequest(http.MethodPost, "/webhooks/1/ativar", nil)
	activateRes := httptest.NewRecorder()
	router.ServeHTTP(activateRes, activateReq)

	require.Equal(t, http.StatusOK, activateRes.Code)
	var activated domainwebhook.WebhookResponse
	require.NoError(t, json.NewDecoder(activateRes.Body).Decode(&activated))
	assert.True(t, activated.Ativo)
}

func TestListWebhooksHandler(t *testing.T) {
	router := newTestRouter()
	createReq := httptest.NewRequest(http.MethodPost, "/webhooks", strings.NewReader(`{"url":"https://example.com/hook"}`))
	router.ServeHTTP(httptest.NewRecorder(), createReq)

	request := httptest.NewRequest(http.MethodGet, "/webhooks", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	var body []domainwebhook.WebhookResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	require.Len(t, body, 1)
	assert.Equal(t, "https://example.com/hook", body[0].URL)
}

func TestCreateJogoEmitsWebhookPayload(t *testing.T) {
	received := make(chan domainwebhook.EventPayload, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload domainwebhook.EventPayload
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		received <- payload
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	router := newTestRouter()
	createWebhook := httptest.NewRequest(http.MethodPost, "/webhooks", strings.NewReader(`{"url":"`+server.URL+`"}`))
	createWebhookRes := httptest.NewRecorder()
	router.ServeHTTP(createWebhookRes, createWebhook)
	require.Equal(t, http.StatusCreated, createWebhookRes.Code)

	createJogo := httptest.NewRequest(http.MethodPost, "/jogos", strings.NewReader(`{
		"nome":"Hades",
		"tipo":"Roguelike",
		"nota":9,
		"review":"Rapido e viciante."
	}`))
	createJogoRes := httptest.NewRecorder()
	router.ServeHTTP(createJogoRes, createJogo)
	require.Equal(t, http.StatusCreated, createJogoRes.Code)

	select {
	case payload := <-received:
		assert.Equal(t, "jogo.criado", payload.Evento)
		assert.Equal(t, 3, payload.ID)
		assert.Equal(t, "Hades", payload.Nome)
	case <-time.After(2 * time.Second):
		t.Fatal("webhook nao recebeu o evento de criacao")
	}
}

func TestDeactivatedWebhookDoesNotReceiveEvents(t *testing.T) {
	received := make(chan struct{}, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received <- struct{}{}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	router := newTestRouter()
	createWebhook := httptest.NewRequest(http.MethodPost, "/webhooks", strings.NewReader(`{"url":"`+server.URL+`"}`))
	router.ServeHTTP(httptest.NewRecorder(), createWebhook)

	deactivate := httptest.NewRequest(http.MethodPost, "/webhooks/1/desativar", nil)
	router.ServeHTTP(httptest.NewRecorder(), deactivate)

	createJogo := httptest.NewRequest(http.MethodPost, "/jogos", strings.NewReader(`{
		"nome":"Hades",
		"tipo":"Roguelike",
		"nota":9,
		"review":"Rapido e viciante."
	}`))
	createJogoRes := httptest.NewRecorder()
	router.ServeHTTP(createJogoRes, createJogo)
	require.Equal(t, http.StatusCreated, createJogoRes.Code)

	select {
	case <-received:
		t.Fatal("webhook desativado nao deveria receber evento")
	case <-time.After(200 * time.Millisecond):
	}
}
