package webhook

import (
	"context"

	domainwebhook "games_api/internal/domain/webhook"
	"games_api/internal/service"
)

type CreateWebhookUseCase struct {
	repo service.WebhookRepository
}

func NewCreateWebhookUseCase(repo service.WebhookRepository) *CreateWebhookUseCase {
	return &CreateWebhookUseCase{repo: repo}
}

func (uc *CreateWebhookUseCase) Execute(ctx context.Context, request domainwebhook.CreateWebhookRequest) (domainwebhook.Webhook, error) {
	if err := request.Validate(); err != nil {
		return domainwebhook.Webhook{}, err
	}

	return uc.repo.Create(ctx, request.ToEntity())
}
