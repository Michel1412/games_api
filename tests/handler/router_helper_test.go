package handler_test

import (
	"net/http"

	"games_api/internal/handler"
	"games_api/internal/service"
	jogousecase "games_api/internal/usecase/jogo"
	userusecase "games_api/internal/usecase/user"
)

func newTestRouter() http.Handler {
	repo := service.NewMemoryJogoRepository()

	return handler.NewRouter(
		userusecase.NewLoginUseCase(),
		jogousecase.NewListJogosUseCase(repo),
		jogousecase.NewGetJogoUseCase(repo),
		jogousecase.NewCreateJogoUseCase(repo),
		jogousecase.NewUpdateJogoUseCase(repo),
		jogousecase.NewDeleteJogoUseCase(repo),
	)
}
