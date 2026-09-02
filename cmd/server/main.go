package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"games_api/internal/bootstrap"
	"games_api/internal/config"
	"games_api/internal/handler"
	"games_api/internal/service"
	jogousecase "games_api/internal/usecase/jogo"
	userusecase "games_api/internal/usecase/user"
	webhookusecase "games_api/internal/usecase/webhook"
)

func main() {
	cfg := config.Load()
	if err := cfg.ValidateProd(); err != nil {
		log.Fatalf("configuracao invalida: %v", err)
	}

	ctx := context.Background()
	jogoRepo, webhookRepo, closeRepos, err := bootstrap.BuildRepositories(ctx, cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer func() {
		if err := closeRepos(); err != nil {
			log.Printf("erro ao fechar repositorio: %v", err)
		}
	}()

	publisher := service.NewWebhookPublisher(webhookRepo, service.NewHTTPWebhookDispatcher())

	router := handler.NewRouter(
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

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("servidor iniciado na porta %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
