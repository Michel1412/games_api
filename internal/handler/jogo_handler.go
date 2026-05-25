package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	domainjogo "games_api/internal/domain/jogo"

	"github.com/go-chi/chi/v5"
)

type JogoHandler struct {
	listUC   ListJogosUseCase
	getUC    GetJogoUseCase
	createUC CreateJogoUseCase
	updateUC UpdateJogoUseCase
	deleteUC DeleteJogoUseCase
}

func NewJogoHandler(
	listUC ListJogosUseCase,
	getUC GetJogoUseCase,
	createUC CreateJogoUseCase,
	updateUC UpdateJogoUseCase,
	deleteUC DeleteJogoUseCase,
) *JogoHandler {
	return &JogoHandler{
		listUC:   listUC,
		getUC:    getUC,
		createUC: createUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
	}
}

func (h *JogoHandler) List(w http.ResponseWriter, r *http.Request) {
	jogos, err := h.listUC.Execute(r.Context())
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, domainjogo.NewJogoResponseList(jogos))
}

func (h *JogoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	jogo, err := h.getUC.Execute(r.Context(), id)
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, domainjogo.NewJogoResponse(jogo))
}

func (h *JogoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request domainjogo.CreateJogoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeBadRequest(w, "json invalido")
		return
	}

	created, err := h.createUC.Execute(r.Context(), request)
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, domainjogo.NewJogoResponse(created))
}

func (h *JogoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var request domainjogo.UpdateJogoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeBadRequest(w, "json invalido")
		return
	}

	updated, err := h.updateUC.Execute(r.Context(), id, request)
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, domainjogo.NewJogoResponse(updated))
}

func (h *JogoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.deleteUC.Execute(r.Context(), id); err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		writeBadRequest(w, "id invalido")
		return 0, false
	}

	return id, true
}
