package jogo

import (
	"context"

	domainwebhook "games_api/internal/domain/webhook"
	"games_api/internal/service"
)

type DeleteJogoUseCase struct {
	repo      service.JogoRepository
	publisher service.JogoEventPublisher
}

func NewDeleteJogoUseCase(repo service.JogoRepository, publisher service.JogoEventPublisher) *DeleteJogoUseCase {
	return &DeleteJogoUseCase{repo: repo, publisher: publisher}
}

func (uc *DeleteJogoUseCase) Execute(ctx context.Context, id int) error {
	jogo, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}

	publishJogoEvent(ctx, uc.publisher, domainwebhook.EventoJogoRemovido, jogo)
	return nil
}
