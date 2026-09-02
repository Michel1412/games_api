package service

import (
	"context"
	"sort"
	"sync"

	"games_api/internal/domain/webhook"
)

type MemoryWebhookRepository struct {
	mu     sync.RWMutex
	nextID int
	items  map[int]webhook.Webhook
}

func NewMemoryWebhookRepository() *MemoryWebhookRepository {
	return &MemoryWebhookRepository{
		nextID: 1,
		items:  make(map[int]webhook.Webhook),
	}
}

func (r *MemoryWebhookRepository) List(_ context.Context) ([]webhook.Webhook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.snapshot(func(item webhook.Webhook) bool { return true }), nil
}

func (r *MemoryWebhookRepository) ListActive(_ context.Context) ([]webhook.Webhook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.snapshot(func(item webhook.Webhook) bool { return item.Ativo }), nil
}

func (r *MemoryWebhookRepository) GetByID(_ context.Context, id int) (webhook.Webhook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	if !ok {
		return webhook.Webhook{}, ErrWebhookNotFound
	}

	return item, nil
}

func (r *MemoryWebhookRepository) Create(_ context.Context, data webhook.Webhook) (webhook.Webhook, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data.ID = r.nextID
	r.items[data.ID] = data
	r.nextID++

	return data, nil
}

func (r *MemoryWebhookRepository) SetAtivo(_ context.Context, id int, ativo bool) (webhook.Webhook, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.items[id]
	if !ok {
		return webhook.Webhook{}, ErrWebhookNotFound
	}

	item.Ativo = ativo
	r.items[id] = item

	return item, nil
}

func (r *MemoryWebhookRepository) snapshot(keep func(webhook.Webhook) bool) []webhook.Webhook {
	ids := make([]int, 0, len(r.items))
	for id := range r.items {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	items := make([]webhook.Webhook, 0, len(ids))
	for _, id := range ids {
		item := r.items[id]
		if keep(item) {
			items = append(items, item)
		}
	}

	return items
}
