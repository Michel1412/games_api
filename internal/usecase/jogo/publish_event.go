package jogo

import (
	"context"

	domainjogo "games_api/internal/domain/jogo"
	domainwebhook "games_api/internal/domain/webhook"
	"games_api/internal/service"
)

func publishJogoEvent(ctx context.Context, publisher service.JogoEventPublisher, evento string, jogo domainjogo.Jogo) {
	if publisher == nil {
		return
	}

	publisher.Publish(ctx, domainwebhook.NewEventPayload(evento, jogo.ID, jogo.Nome))
}
