package handler

import (
	"context"
	"net/http"

	domainjogo "games_api/internal/domain/jogo"
	domainuser "games_api/internal/domain/user"
	domainwebhook "games_api/internal/domain/webhook"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

type ListWebhooksUseCase interface {
	Execute(ctx context.Context) ([]domainwebhook.Webhook, error)
}

type GetWebhookUseCase interface {
	Execute(ctx context.Context, id int) (domainwebhook.Webhook, error)
}

type CreateWebhookUseCase interface {
	Execute(ctx context.Context, request domainwebhook.CreateWebhookRequest) (domainwebhook.Webhook, error)
}

type SetWebhookAtivoUseCase interface {
	Execute(ctx context.Context, id int, ativo bool) (domainwebhook.Webhook, error)
}

func NewRouter(
	loginUC LoginUseCase,
	listUC ListJogosUseCase,
	getUC GetJogoUseCase,
	createUC CreateJogoUseCase,
	updateUC UpdateJogoUseCase,
	deleteUC DeleteJogoUseCase,
	listWebhooksUC ListWebhooksUseCase,
	getWebhookUC GetWebhookUseCase,
	createWebhookUC CreateWebhookUseCase,
	setWebhookAtivoUC SetWebhookAtivoUseCase,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	userHandler := NewUserHandler(loginUC)
	jogoHandler := NewJogoHandler(listUC, getUC, createUC, updateUC, deleteUC)
	webhookHandler := NewWebhookHandler(listWebhooksUC, getWebhookUC, createWebhookUC, setWebhookAtivoUC)

	r.Post("/login", userHandler.Login)
	r.Get("/jogos", jogoHandler.List)
	r.Get("/jogos/{id}", jogoHandler.GetByID)
	r.Post("/jogos", jogoHandler.Create)
	r.Put("/jogos/{id}", jogoHandler.Update)
	r.Delete("/jogos/{id}", jogoHandler.Delete)

	r.Get("/webhooks", webhookHandler.List)
	r.Get("/webhooks/{id}", webhookHandler.GetByID)
	r.Post("/webhooks", webhookHandler.Create)
	r.Post("/webhooks/{id}/ativar", webhookHandler.Activate)
	r.Post("/webhooks/{id}/desativar", webhookHandler.Deactivate)

	return r
}
