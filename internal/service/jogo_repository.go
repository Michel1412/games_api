package service

import (
	"context"
	"errors"

	"games_api/internal/domain/jogo"
)

var ErrJogoNotFound = errors.New("jogo nao encontrado")

type JogoRepository interface {
	List(ctx context.Context) ([]jogo.Jogo, error)
	GetByID(ctx context.Context, id int) (jogo.Jogo, error)
	Create(ctx context.Context, data jogo.Jogo) (jogo.Jogo, error)
	Update(ctx context.Context, id int, data jogo.Jogo) (jogo.Jogo, error)
	Delete(ctx context.Context, id int) error
}
