package jogo

import (
	"context"

	domainjogo "games_api/internal/domain/jogo"
	domainwebhook "games_api/internal/domain/webhook"
	"games_api/internal/service"
)

type UpdateJogoUseCase struct {
	repo      service.JogoRepository
	publisher service.JogoEventPublisher
}

func NewUpdateJogoUseCase(repo service.JogoRepository, publisher service.JogoEventPublisher) *UpdateJogoUseCase {
	return &UpdateJogoUseCase{repo: repo, publisher: publisher}
}

func (uc *UpdateJogoUseCase) Execute(ctx context.Context, id int, request domainjogo.UpdateJogoRequest) (domainjogo.Jogo, error) {
	if err := request.Validate(); err != nil {
		return domainjogo.Jogo{}, err
	}

	updated, err := uc.repo.Update(ctx, id, request.ToEntity(id))
	if err != nil {
		return domainjogo.Jogo{}, err
	}

	publishJogoEvent(ctx, uc.publisher, domainwebhook.EventoJogoAtualizado, updated)
	return updated, nil
}
