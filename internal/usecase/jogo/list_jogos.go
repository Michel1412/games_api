package jogo

import (
	"context"

	domainjogo "games_api/internal/domain/jogo"
	"games_api/internal/service"
)

type ListJogosUseCase struct {
	repo service.JogoRepository
}

func NewListJogosUseCase(repo service.JogoRepository) *ListJogosUseCase {
	return &ListJogosUseCase{repo: repo}
}

func (uc *ListJogosUseCase) Execute(ctx context.Context) ([]domainjogo.Jogo, error) {
	return uc.repo.List(ctx)
}
