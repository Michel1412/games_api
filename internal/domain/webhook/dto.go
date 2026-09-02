package webhook

import (
	"errors"
	"net/url"
	"strings"
)

const (
	EventoJogoCriado     = "jogo.criado"
	EventoJogoAtualizado = "jogo.atualizado"
	EventoJogoRemovido   = "jogo.removido"
)

type CreateWebhookRequest struct {
	URL string `json:"url"`
}

func (r CreateWebhookRequest) Validate() error {
	return validateURL(r.URL)
}

func (r CreateWebhookRequest) ToEntity() Webhook {
	return Webhook{
		URL:   strings.TrimSpace(r.URL),
		Ativo: true,
	}
}

type WebhookResponse struct {
	ID    int    `json:"id"`
	URL   string `json:"url"`
	Ativo bool   `json:"ativo"`
}

func NewWebhookResponse(item Webhook) WebhookResponse {
	return WebhookResponse{
		ID:    item.ID,
		URL:   item.URL,
		Ativo: item.Ativo,
	}
}

func NewWebhookResponseList(items []Webhook) []WebhookResponse {
	response := make([]WebhookResponse, 0, len(items))
	for _, item := range items {
		response = append(response, NewWebhookResponse(item))
	}

	return response
}

type EventPayload struct {
	Evento string `json:"evento"`
	ID     int    `json:"id"`
	Nome   string `json:"nome,omitempty"`
}

func NewEventPayload(evento string, id int, nome string) EventPayload {
	return EventPayload{
		Evento: evento,
		ID:     id,
		Nome:   nome,
	}
}

func validateURL(raw string) error {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return errors.New("url e obrigatoria")
	}

	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil || parsed.Host == "" {
		return errors.New("url invalida")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("url deve usar http ou https")
	}

	return nil
}
