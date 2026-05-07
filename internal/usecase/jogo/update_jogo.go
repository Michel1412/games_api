package jogo

import (
	"context"

	domainjogo "games_api/internal/domain/jogo"
	"games_api/internal/service"
)

type UpdateJogoUseCase struct {
	repo service.JogoRepository
}

func NewUpdateJogoUseCase(repo service.JogoRepository) *UpdateJogoUseCase {
	return &UpdateJogoUseCase{repo: repo}
}

func (uc *UpdateJogoUseCase) Execute(ctx context.Context, id int, request domainjogo.UpdateJogoRequest) (domainjogo.Jogo, error) {
	if err := request.Validate(); err != nil {
		return domainjogo.Jogo{}, err
	}

	return uc.repo.Update(ctx, id, request.ToEntity(id))
}
