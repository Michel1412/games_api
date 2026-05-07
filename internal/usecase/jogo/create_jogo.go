package jogo

import (
	"context"

	domainjogo "games_api/internal/domain/jogo"
	"games_api/internal/service"
)

type CreateJogoUseCase struct {
	repo service.JogoRepository
}

func NewCreateJogoUseCase(repo service.JogoRepository) *CreateJogoUseCase {
	return &CreateJogoUseCase{repo: repo}
}

func (uc *CreateJogoUseCase) Execute(ctx context.Context, request domainjogo.CreateJogoRequest) (domainjogo.Jogo, error) {
	if err := request.Validate(); err != nil {
		return domainjogo.Jogo{}, err
	}

	return uc.repo.Create(ctx, request.ToEntity())
}
