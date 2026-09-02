package jogo

import (
	"context"

	domainjogo "games_api/internal/domain/jogo"
	domainwebhook "games_api/internal/domain/webhook"
	"games_api/internal/service"
)

type CreateJogoUseCase struct {
	repo      service.JogoRepository
	publisher service.JogoEventPublisher
}

func NewCreateJogoUseCase(repo service.JogoRepository, publisher service.JogoEventPublisher) *CreateJogoUseCase {
	return &CreateJogoUseCase{repo: repo, publisher: publisher}
}

func (uc *CreateJogoUseCase) Execute(ctx context.Context, request domainjogo.CreateJogoRequest) (domainjogo.Jogo, error) {
	if err := request.Validate(); err != nil {
		return domainjogo.Jogo{}, err
	}

	created, err := uc.repo.Create(ctx, request.ToEntity())
	if err != nil {
		return domainjogo.Jogo{}, err
	}

	publishJogoEvent(ctx, uc.publisher, domainwebhook.EventoJogoCriado, created)
	return created, nil
}
