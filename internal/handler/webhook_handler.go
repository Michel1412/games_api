package handler

import (
	"encoding/json"
	"net/http"

	domainwebhook "games_api/internal/domain/webhook"
)

type WebhookHandler struct {
	listUC     ListWebhooksUseCase
	getUC      GetWebhookUseCase
	createUC   CreateWebhookUseCase
	setAtivoUC SetWebhookAtivoUseCase
}

func NewWebhookHandler(
	listUC ListWebhooksUseCase,
	getUC GetWebhookUseCase,
	createUC CreateWebhookUseCase,
	setAtivoUC SetWebhookAtivoUseCase,
) *WebhookHandler {
	return &WebhookHandler{
		listUC:     listUC,
		getUC:      getUC,
		createUC:   createUC,
		setAtivoUC: setAtivoUC,
	}
}

func (h *WebhookHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.listUC.Execute(r.Context())
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, domainwebhook.NewWebhookResponseList(items))
}

func (h *WebhookHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	item, err := h.getUC.Execute(r.Context(), id)
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, domainwebhook.NewWebhookResponse(item))
}

func (h *WebhookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request domainwebhook.CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeBadRequest(w, "json invalido")
		return
	}

	created, err := h.createUC.Execute(r.Context(), request)
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, domainwebhook.NewWebhookResponse(created))
}

func (h *WebhookHandler) Activate(w http.ResponseWriter, r *http.Request) {
	h.setAtivo(w, r, true)
}

func (h *WebhookHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	h.setAtivo(w, r, false)
}

func (h *WebhookHandler) setAtivo(w http.ResponseWriter, r *http.Request, ativo bool) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	item, err := h.setAtivoUC.Execute(r.Context(), id, ativo)
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, domainwebhook.NewWebhookResponse(item))
}
