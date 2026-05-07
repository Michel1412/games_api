package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"games_api/internal/config"
	"games_api/internal/handler"
	"games_api/internal/service"
	jogousecase "games_api/internal/usecase/jogo"
	userusecase "games_api/internal/usecase/user"
)

func main() {
	cfg := config.Load()
	if err := cfg.ValidateProd(); err != nil {
		log.Fatalf("configuracao invalida: %v", err)
	}

	ctx := context.Background()
	repo, closeRepo := buildRepository(ctx, cfg)
	defer closeRepo()

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

func buildRepository(ctx context.Context, cfg config.Config) (service.JogoRepository, func()) {
	if !cfg.IsProd() {
		return service.NewMemoryJogoRepository(), func() {}
	}

	if err := cfg.PrepareGoogleCredentialsFile(); err != nil {
		log.Fatalf("erro ao preparar credenciais do firestore: %v", err)
	}

	repo, err := service.NewFirestoreJogoRepository(ctx, cfg.GCPProjectID)
	if err != nil {
		log.Fatalf("erro ao conectar no firestore: %v", err)
	}

	return repo, func() {
		if err := repo.Close(); err != nil {
			log.Printf("erro ao fechar firestore: %v", err)
		}
	}
}
