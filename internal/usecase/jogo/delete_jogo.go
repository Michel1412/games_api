package jogo

import (
	"context"

	"games_api/internal/service"
)

type DeleteJogoUseCase struct {
	repo service.JogoRepository
}

func NewDeleteJogoUseCase(repo service.JogoRepository) *DeleteJogoUseCase {
	return &DeleteJogoUseCase{repo: repo}
}

func (uc *DeleteJogoUseCase) Execute(ctx context.Context, id int) error {
	return uc.repo.Delete(ctx, id)
}
