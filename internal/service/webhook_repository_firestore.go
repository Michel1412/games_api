package service

import (
	"context"
	"strconv"

	"games_api/internal/domain/webhook"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const (
	webhooksCollection = "webhooks"
	webhooksCounterDoc = "webhooks"
)

type FirestoreWebhookRepository struct {
	client *firestore.Client
}

func NewFirestoreWebhookRepository(client *firestore.Client) *FirestoreWebhookRepository {
	return &FirestoreWebhookRepository{client: client}
}

func (r *FirestoreWebhookRepository) List(ctx context.Context) ([]webhook.Webhook, error) {
	iter := r.client.Collection(webhooksCollection).OrderBy("id", firestore.Asc).Documents(ctx)
	defer iter.Stop()

	return collectWebhooks(iter)
}

func (r *FirestoreWebhookRepository) ListActive(ctx context.Context) ([]webhook.Webhook, error) {
	items, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	active := make([]webhook.Webhook, 0, len(items))
	for _, item := range items {
		if item.Ativo {
			active = append(active, item)
		}
	}

	return active, nil
}

func (r *FirestoreWebhookRepository) GetByID(ctx context.Context, id int) (webhook.Webhook, error) {
	doc, err := r.client.Collection(webhooksCollection).Doc(strconv.Itoa(id)).Get(ctx)
	if isNotFound(err) {
		return webhook.Webhook{}, ErrWebhookNotFound
	}
	if err != nil {
		return webhook.Webhook{}, err
	}

	var item webhook.Webhook
	if err := doc.DataTo(&item); err != nil {
		return webhook.Webhook{}, err
	}

	return item, nil
}

func (r *FirestoreWebhookRepository) Create(ctx context.Context, data webhook.Webhook) (webhook.Webhook, error) {
	var created webhook.Webhook

	err := r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		counterRef := r.client.Collection(countersCollection).Doc(webhooksCounterDoc)
		nextID, err := nextFirestoreID(ctx, tx, counterRef)
		if err != nil {
			return err
		}

		data.ID = nextID
		docRef := r.client.Collection(webhooksCollection).Doc(strconv.Itoa(nextID))
		if err := tx.Set(docRef, data); err != nil {
			return err
		}

		if err := tx.Set(counterRef, firestoreCounter{CurrentID: nextID}, firestore.MergeAll); err != nil {
			return err
		}

		created = data
		return nil
	})
	if err != nil {
		return webhook.Webhook{}, err
	}

	return created, nil
}

func (r *FirestoreWebhookRepository) SetAtivo(ctx context.Context, id int, ativo bool) (webhook.Webhook, error) {
	item, err := r.GetByID(ctx, id)
	if err != nil {
		return webhook.Webhook{}, err
	}

	item.Ativo = ativo
	if _, err := r.client.Collection(webhooksCollection).Doc(strconv.Itoa(id)).Set(ctx, item); err != nil {
		return webhook.Webhook{}, err
	}

	return item, nil
}

func collectWebhooks(iter *firestore.DocumentIterator) ([]webhook.Webhook, error) {
	items := make([]webhook.Webhook, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var item webhook.Webhook
		if err := doc.DataTo(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
