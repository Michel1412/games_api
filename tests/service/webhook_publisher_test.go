package service_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"games_api/internal/domain/webhook"
	"games_api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recordingDispatcher struct {
	mu       sync.Mutex
	calls    []string
	payloads [][]byte
}

func (d *recordingDispatcher) Dispatch(url string, body []byte) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.calls = append(d.calls, url)
	d.payloads = append(d.payloads, append([]byte(nil), body...))
}

func TestWebhookPublisherDispatchesOnlyActiveWebhooks(t *testing.T) {
	repo := service.NewMemoryWebhookRepository()
	_, err := repo.Create(context.Background(), webhook.Webhook{URL: "https://active.example/hook", Ativo: true})
	require.NoError(t, err)
	createdInactive, err := repo.Create(context.Background(), webhook.Webhook{URL: "https://inactive.example/hook", Ativo: true})
	require.NoError(t, err)
	_, err = repo.SetAtivo(context.Background(), createdInactive.ID, false)
	require.NoError(t, err)

	dispatcher := &recordingDispatcher{}
	publisher := service.NewWebhookPublisher(repo, dispatcher)

	publisher.Publish(context.Background(), webhook.NewEventPayload(webhook.EventoJogoCriado, 3, "Elden Ring"))

	require.Eventually(t, func() bool {
		dispatcher.mu.Lock()
		defer dispatcher.mu.Unlock()
		return len(dispatcher.calls) == 1
	}, time.Second, 10*time.Millisecond)

	dispatcher.mu.Lock()
	defer dispatcher.mu.Unlock()
	assert.Equal(t, []string{"https://active.example/hook"}, dispatcher.calls)

	var payload webhook.EventPayload
	require.NoError(t, json.Unmarshal(dispatcher.payloads[0], &payload))
	assert.Equal(t, "jogo.criado", payload.Evento)
	assert.Equal(t, 3, payload.ID)
	assert.Equal(t, "Elden Ring", payload.Nome)
}

func TestHTTPWebhookDispatcherPostsJSON(t *testing.T) {
	received := make(chan []byte, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		received <- body
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	dispatcher := service.NewHTTPWebhookDispatcher()
	payload, err := json.Marshal(webhook.NewEventPayload(webhook.EventoJogoAtualizado, 1, "Zelda"))
	require.NoError(t, err)

	dispatcher.Dispatch(server.URL, payload)

	select {
	case body := <-received:
		var event webhook.EventPayload
		require.NoError(t, json.Unmarshal(body, &event))
		assert.Equal(t, "jogo.atualizado", event.Evento)
		assert.Equal(t, 1, event.ID)
	case <-time.After(time.Second):
		t.Fatal("dispatcher nao enviou o payload")
	}
}
