package webhook_test

import (
	"context"
	"testing"

	domainwebhook "games_api/internal/domain/webhook"
	"games_api/internal/service"
	webhookusecase "games_api/internal/usecase/webhook"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateWebhookUseCaseValidatesURL(t *testing.T) {
	uc := webhookusecase.NewCreateWebhookUseCase(service.NewMemoryWebhookRepository())

	_, err := uc.Execute(context.Background(), domainwebhook.CreateWebhookRequest{URL: "ftp://x"})

	require.Error(t, err)
	assert.Equal(t, "url deve usar http ou https", err.Error())
}

func TestCreateWebhookUseCaseCreatesActiveWebhook(t *testing.T) {
	repo := service.NewMemoryWebhookRepository()
	uc := webhookusecase.NewCreateWebhookUseCase(repo)

	created, err := uc.Execute(context.Background(), domainwebhook.CreateWebhookRequest{
		URL: "https://example.com/hook",
	})

	require.NoError(t, err)
	assert.Equal(t, 1, created.ID)
	assert.True(t, created.Ativo)
}

func TestSetWebhookAtivoUseCase(t *testing.T) {
	repo := service.NewMemoryWebhookRepository()
	created, err := repo.Create(context.Background(), domainwebhook.Webhook{URL: "https://example.com/hook", Ativo: true})
	require.NoError(t, err)

	uc := webhookusecase.NewSetWebhookAtivoUseCase(repo)
	updated, err := uc.Execute(context.Background(), created.ID, false)

	require.NoError(t, err)
	assert.False(t, updated.Ativo)
}

func TestGetWebhookUseCaseNotFound(t *testing.T) {
	uc := webhookusecase.NewGetWebhookUseCase(service.NewMemoryWebhookRepository())

	_, err := uc.Execute(context.Background(), 9)

	assert.ErrorIs(t, err, service.ErrWebhookNotFound)
}

func TestListWebhooksUseCase(t *testing.T) {
	repo := service.NewMemoryWebhookRepository()
	_, err := repo.Create(context.Background(), domainwebhook.Webhook{URL: "https://example.com/hook", Ativo: true})
	require.NoError(t, err)

	uc := webhookusecase.NewListWebhooksUseCase(repo)
	items, err := uc.Execute(context.Background())

	require.NoError(t, err)
	require.Len(t, items, 1)
}
