package service

import (
	"context"

	"cloud.google.com/go/firestore"
)

func NewFirestoreRepositories(ctx context.Context, projectID string) (*FirestoreJogoRepository, *FirestoreWebhookRepository, func() error, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, nil, nil, err
	}

	return NewFirestoreJogoRepositoryWithClient(client), NewFirestoreWebhookRepository(client), client.Close, nil
}

func NewFirestoreJogoRepositoryWithClient(client *firestore.Client) *FirestoreJogoRepository {
	return &FirestoreJogoRepository{client: client}
}

func nextFirestoreID(ctx context.Context, tx *firestore.Transaction, counterRef *firestore.DocumentRef) (int, error) {
	doc, err := tx.Get(counterRef)
	if isNotFound(err) {
		return 1, nil
	}
	if err != nil {
		return 0, err
	}

	var counter firestoreCounter
	if err := doc.DataTo(&counter); err != nil {
		return 0, err
	}

	return counter.CurrentID + 1, nil
}
