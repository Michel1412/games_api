package handler_test

import (
	"net/http"

	"games_api/internal/handler"
	"games_api/internal/service"
	jogousecase "games_api/internal/usecase/jogo"
	userusecase "games_api/internal/usecase/user"
	webhookusecase "games_api/internal/usecase/webhook"
)

func newTestRouter() http.Handler {
	jogoRepo := service.NewMemoryJogoRepository()
	webhookRepo := service.NewMemoryWebhookRepository()
	publisher := service.NewWebhookPublisher(webhookRepo, service.NewHTTPWebhookDispatcher())

	return handler.NewRouter(
		userusecase.NewLoginUseCase(),
		jogousecase.NewListJogosUseCase(jogoRepo),
		jogousecase.NewGetJogoUseCase(jogoRepo),
		jogousecase.NewCreateJogoUseCase(jogoRepo, publisher),
		jogousecase.NewUpdateJogoUseCase(jogoRepo, publisher),
		jogousecase.NewDeleteJogoUseCase(jogoRepo, publisher),
		webhookusecase.NewListWebhooksUseCase(webhookRepo),
		webhookusecase.NewGetWebhookUseCase(webhookRepo),
		webhookusecase.NewCreateWebhookUseCase(webhookRepo),
		webhookusecase.NewSetWebhookAtivoUseCase(webhookRepo),
	)
}
