package jogo

import (
	"errors"
	"strings"
)

type CreateJogoRequest struct {
	Nome   string `json:"nome"`
	Tipo   string `json:"tipo"`
	Nota   int    `json:"nota"`
	Review string `json:"review"`
}

func (r CreateJogoRequest) Validate() error {
	return validateFields(r.Nome, r.Tipo, r.Nota, r.Review)
}

func (r CreateJogoRequest) ToEntity() Jogo {
	return Jogo{
		Nome:   strings.TrimSpace(r.Nome),
		Tipo:   strings.TrimSpace(r.Tipo),
		Nota:   r.Nota,
		Review: strings.TrimSpace(r.Review),
	}
}

type UpdateJogoRequest struct {
	Nome   string `json:"nome"`
	Tipo   string `json:"tipo"`
	Nota   int    `json:"nota"`
	Review string `json:"review"`
}

func (r UpdateJogoRequest) Validate() error {
	return validateFields(r.Nome, r.Tipo, r.Nota, r.Review)
}

func (r UpdateJogoRequest) ToEntity(id int) Jogo {
	return Jogo{
		ID:     id,
		Nome:   strings.TrimSpace(r.Nome),
		Tipo:   strings.TrimSpace(r.Tipo),
		Nota:   r.Nota,
		Review: strings.TrimSpace(r.Review),
	}
}

type JogoResponse struct {
	ID     int    `json:"id"`
	Nome   string `json:"nome"`
	Tipo   string `json:"tipo"`
	Nota   int    `json:"nota"`
	Review string `json:"review"`
}

func NewJogoResponse(j Jogo) JogoResponse {
	return JogoResponse{
		ID:     j.ID,
		Nome:   j.Nome,
		Tipo:   j.Tipo,
		Nota:   j.Nota,
		Review: j.Review,
	}
}

func NewJogoResponseList(jogos []Jogo) []JogoResponse {
	response := make([]JogoResponse, 0, len(jogos))
	for _, item := range jogos {
		response = append(response, NewJogoResponse(item))
	}

	return response
}

func validateFields(nome, tipo string, nota int, review string) error {
	if strings.TrimSpace(nome) == "" {
		return errors.New("nome e obrigatorio")
	}

	if strings.TrimSpace(tipo) == "" {
		return errors.New("tipo e obrigatorio")
	}

	if nota < 1 || nota > 10 {
		return errors.New("nota deve estar entre 1 e 10")
	}

	if strings.TrimSpace(review) == "" {
		return errors.New("review e obrigatoria")
	}

	return nil
}
