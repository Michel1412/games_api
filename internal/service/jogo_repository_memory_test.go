package service

import (
	"context"
	"testing"

	"games_api/internal/domain/jogo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryJogoRepositoryListReturnsInitialJogos(t *testing.T) {
	repo := NewMemoryJogoRepository()

	jogos, err := repo.List(context.Background())

	require.NoError(t, err)
	require.Len(t, jogos, 2)
	assert.Equal(t, 1, jogos[0].ID)
	assert.Equal(t, "The Legend of Zelda", jogos[0].Nome)
	assert.Equal(t, 2, jogos[1].ID)
}

func TestMemoryJogoRepositoryCreateUsesSequentialID(t *testing.T) {
	repo := NewMemoryJogoRepository()

	created, err := repo.Create(context.Background(), jogo.Jogo{
		Nome:   "Elden Ring",
		Tipo:   "RPG",
		Nota:   9,
		Review: "Desafiador e visualmente impecavel.",
	})

	require.NoError(t, err)
	assert.Equal(t, 3, created.ID)
}

func TestMemoryJogoRepositoryCRUD(t *testing.T) {
	repo := NewMemoryJogoRepository()

	created, err := repo.Create(context.Background(), jogo.Jogo{
		Nome:   "Hades",
		Tipo:   "Roguelike",
		Nota:   9,
		Review: "Rapido e viciante.",
	})
	require.NoError(t, err)

	found, err := repo.GetByID(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Hades", found.Nome)

	updated, err := repo.Update(context.Background(), created.ID, jogo.Jogo{
		Nome:   "Hades II",
		Tipo:   "Roguelike",
		Nota:   10,
		Review: "Ainda melhor.",
	})
	require.NoError(t, err)
	assert.Equal(t, "Hades II", updated.Nome)

	err = repo.Delete(context.Background(), created.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(context.Background(), created.ID)
	assert.ErrorIs(t, err, ErrJogoNotFound)
}
