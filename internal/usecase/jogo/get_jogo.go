package jogo

import (
	"context"

	domainjogo "games_api/internal/domain/jogo"
	"games_api/internal/service"
)

type GetJogoUseCase struct {
	repo service.JogoRepository
}

func NewGetJogoUseCase(repo service.JogoRepository) *GetJogoUseCase {
	return &GetJogoUseCase{repo: repo}
}

func (uc *GetJogoUseCase) Execute(ctx context.Context, id int) (domainjogo.Jogo, error) {
	return uc.repo.GetByID(ctx, id)
}
