package webhook

import (
	"context"

	domainwebhook "games_api/internal/domain/webhook"
	"games_api/internal/service"
)

type GetWebhookUseCase struct {
	repo service.WebhookRepository
}

func NewGetWebhookUseCase(repo service.WebhookRepository) *GetWebhookUseCase {
	return &GetWebhookUseCase{repo: repo}
}

func (uc *GetWebhookUseCase) Execute(ctx context.Context, id int) (domainwebhook.Webhook, error) {
	return uc.repo.GetByID(ctx, id)
}
