package bootstrap

import (
	"context"
	"fmt"

	"games_api/internal/config"
	"games_api/internal/service"
)

func BuildRepository(ctx context.Context, cfg config.Config) (service.JogoRepository, func() error, error) {
	if !cfg.IsProd() {
		return service.NewMemoryJogoRepository(), func() error { return nil }, nil
	}

	if err := cfg.PrepareGoogleCredentialsFile(); err != nil {
		return nil, nil, fmt.Errorf("erro ao preparar credenciais do firestore: %w", err)
	}

	repo, err := service.NewFirestoreJogoRepository(ctx, cfg.GCPProjectID)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao conectar no firestore: %w", err)
	}

	return repo, repo.Close, nil
}
