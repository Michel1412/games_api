package webhook

import (
	"context"

	domainwebhook "games_api/internal/domain/webhook"
	"games_api/internal/service"
)

type SetWebhookAtivoUseCase struct {
	repo service.WebhookRepository
}

func NewSetWebhookAtivoUseCase(repo service.WebhookRepository) *SetWebhookAtivoUseCase {
	return &SetWebhookAtivoUseCase{repo: repo}
}

func (uc *SetWebhookAtivoUseCase) Execute(ctx context.Context, id int, ativo bool) (domainwebhook.Webhook, error) {
	return uc.repo.SetAtivo(ctx, id, ativo)
}
