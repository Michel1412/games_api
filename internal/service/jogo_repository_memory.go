package service

import (
	"context"
	"sort"
	"sync"

	"games_api/internal/domain/jogo"
)

type MemoryJogoRepository struct {
	mu     sync.RWMutex
	nextID int
	items  map[int]jogo.Jogo
}

func NewMemoryJogoRepository() *MemoryJogoRepository {
	initial := []jogo.Jogo{
		{
			ID:     1,
			Nome:   "The Legend of Zelda",
			Tipo:   "Aventura",
			Nota:   10,
			Review: "Um classico absoluto.",
		},
		{
			ID:     2,
			Nome:   "FIFA 23",
			Tipo:   "Esporte",
			Nota:   7,
			Review: "Bom para jogar com amigos.",
		},
	}

	items := make(map[int]jogo.Jogo, len(initial))
	for _, item := range initial {
		items[item.ID] = item
	}

	return &MemoryJogoRepository{
		nextID: 3,
		items:  items,
	}
}

func (r *MemoryJogoRepository) List(_ context.Context) ([]jogo.Jogo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]int, 0, len(r.items))
	for id := range r.items {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	jogos := make([]jogo.Jogo, 0, len(ids))
	for _, id := range ids {
		jogos = append(jogos, r.items[id])
	}

	return jogos, nil
}

func (r *MemoryJogoRepository) GetByID(_ context.Context, id int) (jogo.Jogo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	if !ok {
		return jogo.Jogo{}, ErrJogoNotFound
	}

	return item, nil
}

func (r *MemoryJogoRepository) Create(_ context.Context, data jogo.Jogo) (jogo.Jogo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data.ID = r.nextID
	r.items[data.ID] = data
	r.nextID++

	return data, nil
}

func (r *MemoryJogoRepository) Update(_ context.Context, id int, data jogo.Jogo) (jogo.Jogo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return jogo.Jogo{}, ErrJogoNotFound
	}

	data.ID = id
	r.items[id] = data

	return data, nil
}

func (r *MemoryJogoRepository) Delete(_ context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return ErrJogoNotFound
	}

	delete(r.items, id)
	return nil
}
