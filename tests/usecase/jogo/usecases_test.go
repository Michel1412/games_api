package jogo_test

import (
	"context"
	"testing"

	domainjogo "games_api/internal/domain/jogo"
	"games_api/internal/service"
	jogousecase "games_api/internal/usecase/jogo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeJogoRepository struct {
	items        []domainjogo.Jogo
	created      domainjogo.Jogo
	updated      domainjogo.Jogo
	err          error
	createCalled bool
	updateCalled bool
	deleteCalled bool
	receivedID   int
	receivedJogo domainjogo.Jogo
}

func (r *fakeJogoRepository) List(_ context.Context) ([]domainjogo.Jogo, error) {
	return r.items, r.err
}

func (r *fakeJogoRepository) GetByID(_ context.Context, id int) (domainjogo.Jogo, error) {
	r.receivedID = id
	if r.err != nil {
		return domainjogo.Jogo{}, r.err
	}
	return r.items[0], nil
}

func (r *fakeJogoRepository) Create(_ context.Context, data domainjogo.Jogo) (domainjogo.Jogo, error) {
	r.createCalled = true
	r.receivedJogo = data
	if r.err != nil {
		return domainjogo.Jogo{}, r.err
	}
	return r.created, nil
}

func (r *fakeJogoRepository) Update(_ context.Context, id int, data domainjogo.Jogo) (domainjogo.Jogo, error) {
	r.updateCalled = true
	r.receivedID = id
	r.receivedJogo = data
	if r.err != nil {
		return domainjogo.Jogo{}, r.err
	}
	return r.updated, nil
}

func (r *fakeJogoRepository) Delete(_ context.Context, id int) error {
	r.deleteCalled = true
	r.receivedID = id
	return r.err
}

func TestListJogosUseCase(t *testing.T) {
	repo := &fakeJogoRepository{items: []domainjogo.Jogo{{ID: 1, Nome: "Zelda"}}}
	uc := jogousecase.NewListJogosUseCase(repo)

	result, err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Zelda", result[0].Nome)
}

func TestGetJogoUseCase(t *testing.T) {
	repo := &fakeJogoRepository{items: []domainjogo.Jogo{{ID: 2, Nome: "FIFA 23"}}}
	uc := jogousecase.NewGetJogoUseCase(repo)

	result, err := uc.Execute(context.Background(), 2)

	require.NoError(t, err)
	assert.Equal(t, 2, repo.receivedID)
	assert.Equal(t, "FIFA 23", result.Nome)
}

func TestCreateJogoUseCaseValidatesRequest(t *testing.T) {
	repo := &fakeJogoRepository{}
	uc := jogousecase.NewCreateJogoUseCase(repo)

	_, err := uc.Execute(context.Background(), domainjogo.CreateJogoRequest{Nota: 11})

	require.Error(t, err)
	assert.False(t, repo.createCalled)
}

func TestCreateJogoUseCaseCreatesJogo(t *testing.T) {
	repo := &fakeJogoRepository{created: domainjogo.Jogo{ID: 3, Nome: "Elden Ring", Tipo: "RPG", Nota: 9, Review: "Otimo"}}
	uc := jogousecase.NewCreateJogoUseCase(repo)

	result, err := uc.Execute(context.Background(), domainjogo.CreateJogoRequest{
		Nome: " Elden Ring ", Tipo: "RPG", Nota: 9, Review: "Otimo",
	})

	require.NoError(t, err)
	assert.True(t, repo.createCalled)
	assert.Equal(t, "Elden Ring", repo.receivedJogo.Nome)
	assert.Equal(t, 3, result.ID)
}

func TestUpdateJogoUseCaseValidatesRequest(t *testing.T) {
	repo := &fakeJogoRepository{}
	uc := jogousecase.NewUpdateJogoUseCase(repo)

	_, err := uc.Execute(context.Background(), 1, domainjogo.UpdateJogoRequest{Nota: 0})

	require.Error(t, err)
	assert.False(t, repo.updateCalled)
}

func TestUpdateJogoUseCaseReturnsRepositoryError(t *testing.T) {
	repo := &fakeJogoRepository{err: service.ErrJogoNotFound}
	uc := jogousecase.NewUpdateJogoUseCase(repo)

	_, err := uc.Execute(context.Background(), 99, domainjogo.UpdateJogoRequest{
		Nome: "Halo", Tipo: "FPS", Nota: 8, Review: "Classico",
	})

	assert.ErrorIs(t, err, service.ErrJogoNotFound)
	assert.True(t, repo.updateCalled)
}

func TestDeleteJogoUseCase(t *testing.T) {
	repo := &fakeJogoRepository{}
	uc := jogousecase.NewDeleteJogoUseCase(repo)

	err := uc.Execute(context.Background(), 2)

	require.NoError(t, err)
	assert.True(t, repo.deleteCalled)
	assert.Equal(t, 2, repo.receivedID)
}
