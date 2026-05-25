package bootstrap_test

import (
	"context"
	"testing"

	"games_api/internal/bootstrap"
	"games_api/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildRepositoryUsesMemoryWhenNotProd(t *testing.T) {
	repo, closeRepo, err := bootstrap.BuildRepository(context.Background(), config.Config{Env: "dev"})
	require.NoError(t, err)
	require.NotNil(t, repo)
	require.NotNil(t, closeRepo)

	jogos, err := repo.List(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, jogos)

	assert.NoError(t, closeRepo())
}
