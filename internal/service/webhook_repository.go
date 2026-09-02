package service

import (
	"context"
	"errors"

	"games_api/internal/domain/webhook"
)

var ErrWebhookNotFound = errors.New("webhook nao encontrado")

type WebhookRepository interface {
	List(ctx context.Context) ([]webhook.Webhook, error)
	ListActive(ctx context.Context) ([]webhook.Webhook, error)
	GetByID(ctx context.Context, id int) (webhook.Webhook, error)
	Create(ctx context.Context, data webhook.Webhook) (webhook.Webhook, error)
	SetAtivo(ctx context.Context, id int, ativo bool) (webhook.Webhook, error)
}

type JogoEventPublisher interface {
	Publish(ctx context.Context, payload webhook.EventPayload)
}
