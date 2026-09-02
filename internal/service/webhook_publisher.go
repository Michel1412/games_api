package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"games_api/internal/domain/webhook"
)

type WebhookDispatcher interface {
	Dispatch(url string, body []byte)
}

type HTTPWebhookDispatcher struct {
	client *http.Client
}

func NewHTTPWebhookDispatcher() *HTTPWebhookDispatcher {
	return &HTTPWebhookDispatcher{
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (d *HTTPWebhookDispatcher) Dispatch(targetURL string, body []byte) {
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(body))
	if err != nil {
		log.Printf("erro ao criar requisicao de webhook: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		log.Printf("erro ao enviar webhook para %s: %v", targetURL, err)
		return
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= http.StatusMultipleChoices {
		log.Printf("webhook %s retornou status %d", targetURL, resp.StatusCode)
	}
}

type WebhookPublisher struct {
	repo       WebhookRepository
	dispatcher WebhookDispatcher
}

func NewWebhookPublisher(repo WebhookRepository, dispatcher WebhookDispatcher) *WebhookPublisher {
	return &WebhookPublisher{repo: repo, dispatcher: dispatcher}
}

func (p *WebhookPublisher) Publish(ctx context.Context, payload webhook.EventPayload) {
	if p == nil || p.repo == nil || p.dispatcher == nil {
		return
	}

	hooks, err := p.repo.ListActive(ctx)
	if err != nil {
		log.Printf("erro ao listar webhooks ativos: %v", err)
		return
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("erro ao serializar payload do webhook: %v", err)
		return
	}

	for _, hook := range hooks {
		hook := hook
		go p.dispatcher.Dispatch(hook.URL, body)
	}
}
