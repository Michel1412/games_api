package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"games_api/internal/bootstrap"
	"games_api/internal/config"
	"games_api/internal/handler"
	jogousecase "games_api/internal/usecase/jogo"
	userusecase "games_api/internal/usecase/user"
)

func main() {
	cfg := config.Load()
	if err := cfg.ValidateProd(); err != nil {
		log.Fatalf("configuracao invalida: %v", err)
	}

	ctx := context.Background()
	repo, closeRepo, err := bootstrap.BuildRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer func() {
		if err := closeRepo(); err != nil {
			log.Printf("erro ao fechar repositorio: %v", err)
		}
	}()

	router := handler.NewRouter(
		userusecase.NewLoginUseCase(),
		jogousecase.NewListJogosUseCase(repo),
		jogousecase.NewGetJogoUseCase(repo),
		jogousecase.NewCreateJogoUseCase(repo),
		jogousecase.NewUpdateJogoUseCase(repo),
		jogousecase.NewDeleteJogoUseCase(repo),
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
