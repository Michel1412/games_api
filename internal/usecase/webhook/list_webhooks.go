package webhook

import (
	"context"

	domainwebhook "games_api/internal/domain/webhook"
	"games_api/internal/service"
)

type ListWebhooksUseCase struct {
	repo service.WebhookRepository
}

func NewListWebhooksUseCase(repo service.WebhookRepository) *ListWebhooksUseCase {
	return &ListWebhooksUseCase{repo: repo}
}

func (uc *ListWebhooksUseCase) Execute(ctx context.Context) ([]domainwebhook.Webhook, error) {
	return uc.repo.List(ctx)
}
