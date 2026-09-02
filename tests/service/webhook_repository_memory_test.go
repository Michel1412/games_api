package service_test

import (
	"context"
	"testing"

	"games_api/internal/domain/webhook"
	"games_api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryWebhookRepositoryCreateAndList(t *testing.T) {
	repo := service.NewMemoryWebhookRepository()

	created, err := repo.Create(context.Background(), webhook.Webhook{
		URL:   "https://example.com/hook",
		Ativo: true,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, created.ID)

	items, err := repo.List(context.Background())
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "https://example.com/hook", items[0].URL)
}

func TestMemoryWebhookRepositorySetAtivoAndListActive(t *testing.T) {
	repo := service.NewMemoryWebhookRepository()
	created, err := repo.Create(context.Background(), webhook.Webhook{
		URL:   "https://example.com/hook",
		Ativo: true,
	})
	require.NoError(t, err)

	updated, err := repo.SetAtivo(context.Background(), created.ID, false)
	require.NoError(t, err)
	assert.False(t, updated.Ativo)

	active, err := repo.ListActive(context.Background())
	require.NoError(t, err)
	assert.Empty(t, active)

	_, err = repo.SetAtivo(context.Background(), created.ID, true)
	require.NoError(t, err)

	active, err = repo.ListActive(context.Background())
	require.NoError(t, err)
	require.Len(t, active, 1)
}

func TestMemoryWebhookRepositoryGetByIDNotFound(t *testing.T) {
	repo := service.NewMemoryWebhookRepository()

	_, err := repo.GetByID(context.Background(), 99)

	assert.ErrorIs(t, err, service.ErrWebhookNotFound)
}
