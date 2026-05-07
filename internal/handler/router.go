package handler

import (
	"context"
	"net/http"

	domainjogo "games_api/internal/domain/jogo"
	domainuser "games_api/internal/domain/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type LoginUseCase interface {
	Execute(request domainuser.LoginRequest) (domainuser.LoginResponse, error)
}

type ListJogosUseCase interface {
	Execute(ctx context.Context) ([]domainjogo.Jogo, error)
}

type GetJogoUseCase interface {
	Execute(ctx context.Context, id int) (domainjogo.Jogo, error)
}

type CreateJogoUseCase interface {
	Execute(ctx context.Context, request domainjogo.CreateJogoRequest) (domainjogo.Jogo, error)
}

type UpdateJogoUseCase interface {
	Execute(ctx context.Context, id int, request domainjogo.UpdateJogoRequest) (domainjogo.Jogo, error)
}

type DeleteJogoUseCase interface {
	Execute(ctx context.Context, id int) error
}

func NewRouter(
	loginUC LoginUseCase,
	listUC ListJogosUseCase,
	getUC GetJogoUseCase,
	createUC CreateJogoUseCase,
	updateUC UpdateJogoUseCase,
	deleteUC DeleteJogoUseCase,
) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	userHandler := NewUserHandler(loginUC)
	jogoHandler := NewJogoHandler(listUC, getUC, createUC, updateUC, deleteUC)

	r.Post("/login", userHandler.Login)
	r.Get("/jogos", jogoHandler.List)
	r.Get("/jogos/{id}", jogoHandler.GetByID)
	r.Post("/jogos", jogoHandler.Create)
	r.Put("/jogos/{id}", jogoHandler.Update)
	r.Delete("/jogos/{id}", jogoHandler.Delete)

	return r
}
