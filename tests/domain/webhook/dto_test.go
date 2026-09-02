package webhook_test

import (
	"testing"

	domainwebhook "games_api/internal/domain/webhook"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateWebhookRequestValidate(t *testing.T) {
	cases := []struct {
		name    string
		request domainwebhook.CreateWebhookRequest
		wantErr string
	}{
		{
			name:    "https valido",
			request: domainwebhook.CreateWebhookRequest{URL: "https://example.com/hook"},
		},
		{
			name:    "http valido",
			request: domainwebhook.CreateWebhookRequest{URL: "http://localhost:9090/hook"},
		},
		{
			name:    "vazio",
			request: domainwebhook.CreateWebhookRequest{URL: "   "},
			wantErr: "url e obrigatoria",
		},
		{
			name:    "invalida",
			request: domainwebhook.CreateWebhookRequest{URL: "not-a-url"},
			wantErr: "url invalida",
		},
		{
			name:    "esquema invalido",
			request: domainwebhook.CreateWebhookRequest{URL: "ftp://example.com/hook"},
			wantErr: "url deve usar http ou https",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.request.Validate()
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}

			require.Error(t, err)
			assert.Equal(t, tc.wantErr, err.Error())
		})
	}
}

func TestCreateWebhookRequestToEntityStartsActive(t *testing.T) {
	request := domainwebhook.CreateWebhookRequest{URL: " https://example.com/hook "}

	entity := request.ToEntity()

	assert.Equal(t, "https://example.com/hook", entity.URL)
	assert.True(t, entity.Ativo)
}

func TestNewWebhookResponseList(t *testing.T) {
	items := []domainwebhook.Webhook{
		{ID: 1, URL: "https://a.example", Ativo: true},
		{ID: 2, URL: "https://b.example", Ativo: false},
	}

	response := domainwebhook.NewWebhookResponseList(items)

	require.Len(t, response, 2)
	assert.Equal(t, 1, response[0].ID)
	assert.False(t, response[1].Ativo)
}

func TestNewEventPayload(t *testing.T) {
	payload := domainwebhook.NewEventPayload(domainwebhook.EventoJogoCriado, 3, "Elden Ring")

	assert.Equal(t, "jogo.criado", payload.Evento)
	assert.Equal(t, 3, payload.ID)
	assert.Equal(t, "Elden Ring", payload.Nome)
}
