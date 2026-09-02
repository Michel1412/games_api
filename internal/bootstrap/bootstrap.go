package bootstrap

import (
	"context"
	"fmt"

	"games_api/internal/config"
	"games_api/internal/service"
)

func BuildRepositories(ctx context.Context, cfg config.Config) (service.JogoRepository, service.WebhookRepository, func() error, error) {
	if !cfg.IsProd() {
		return service.NewMemoryJogoRepository(), service.NewMemoryWebhookRepository(), func() error { return nil }, nil
	}

	if err := cfg.PrepareGoogleCredentialsFile(); err != nil {
		return nil, nil, nil, fmt.Errorf("erro ao preparar credenciais do firestore: %w", err)
	}

	jogoRepo, webhookRepo, closeRepos, err := service.NewFirestoreRepositories(ctx, cfg.GCPProjectID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("erro ao conectar no firestore: %w", err)
	}

	return jogoRepo, webhookRepo, closeRepos, nil
}
