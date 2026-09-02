package service

import (
	"context"
	"strconv"

	"games_api/internal/domain/jogo"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	jogosCollection    = "jogos"
	countersCollection = "counters"
	jogosCounterDoc    = "jogos"
)

type FirestoreJogoRepository struct {
	client *firestore.Client
}

type firestoreCounter struct {
	CurrentID int `firestore:"current_id"`
}

func NewFirestoreJogoRepository(ctx context.Context, projectID string) (*FirestoreJogoRepository, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return NewFirestoreJogoRepositoryWithClient(client), nil
}

func (r *FirestoreJogoRepository) Close() error {
	return r.client.Close()
}

func (r *FirestoreJogoRepository) List(ctx context.Context) ([]jogo.Jogo, error) {
	iter := r.client.Collection(jogosCollection).OrderBy("id", firestore.Asc).Documents(ctx)
	defer iter.Stop()

	jogos := make([]jogo.Jogo, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var item jogo.Jogo
		if err := doc.DataTo(&item); err != nil {
			return nil, err
		}
		jogos = append(jogos, item)
	}

	return jogos, nil
}

func (r *FirestoreJogoRepository) GetByID(ctx context.Context, id int) (jogo.Jogo, error) {
	doc, err := r.client.Collection(jogosCollection).Doc(strconv.Itoa(id)).Get(ctx)
	if isNotFound(err) {
		return jogo.Jogo{}, ErrJogoNotFound
	}
	if err != nil {
		return jogo.Jogo{}, err
	}

	var item jogo.Jogo
	if err := doc.DataTo(&item); err != nil {
		return jogo.Jogo{}, err
	}

	return item, nil
}

func (r *FirestoreJogoRepository) Create(ctx context.Context, data jogo.Jogo) (jogo.Jogo, error) {
	var created jogo.Jogo

	err := r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		counterRef := r.client.Collection(countersCollection).Doc(jogosCounterDoc)
		nextID, err := nextFirestoreID(ctx, tx, counterRef)
		if err != nil {
			return err
		}

		data.ID = nextID
		docRef := r.client.Collection(jogosCollection).Doc(strconv.Itoa(nextID))
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
		return jogo.Jogo{}, err
	}

	return created, nil
}

func (r *FirestoreJogoRepository) Update(ctx context.Context, id int, data jogo.Jogo) (jogo.Jogo, error) {
	docRef := r.client.Collection(jogosCollection).Doc(strconv.Itoa(id))
	if _, err := docRef.Get(ctx); isNotFound(err) {
		return jogo.Jogo{}, ErrJogoNotFound
	} else if err != nil {
		return jogo.Jogo{}, err
	}

	data.ID = id
	if _, err := docRef.Set(ctx, data); err != nil {
		return jogo.Jogo{}, err
	}

	return data, nil
}

func (r *FirestoreJogoRepository) Delete(ctx context.Context, id int) error {
	docRef := r.client.Collection(jogosCollection).Doc(strconv.Itoa(id))
	if _, err := docRef.Get(ctx); isNotFound(err) {
		return ErrJogoNotFound
	} else if err != nil {
		return err
	}

	_, err := docRef.Delete(ctx)
	return err
}

func isNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}
