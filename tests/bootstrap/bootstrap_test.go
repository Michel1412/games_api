package bootstrap_test

import (
	"context"
	"testing"

	"games_api/internal/bootstrap"
	"games_api/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildRepositoriesUsesMemoryWhenNotProd(t *testing.T) {
	jogoRepo, webhookRepo, closeRepos, err := bootstrap.BuildRepositories(context.Background(), config.Config{Env: "dev"})
	require.NoError(t, err)
	require.NotNil(t, jogoRepo)
	require.NotNil(t, webhookRepo)
	require.NotNil(t, closeRepos)

	jogos, err := jogoRepo.List(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, jogos)

	webhooks, err := webhookRepo.List(context.Background())
	require.NoError(t, err)
	assert.Empty(t, webhooks)

	assert.NoError(t, closeRepos())
}
