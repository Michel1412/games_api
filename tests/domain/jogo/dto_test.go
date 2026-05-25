package jogo_test

import (
	"testing"

	domainjogo "games_api/internal/domain/jogo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateJogoRequestValidate(t *testing.T) {
	cases := []struct {
		name    string
		request domainjogo.CreateJogoRequest
		wantErr string
	}{
		{
			name:    "valido",
			request: domainjogo.CreateJogoRequest{Nome: "Elden Ring", Tipo: "RPG", Nota: 9, Review: "Otimo"},
		},
		{
			name:    "nome vazio",
			request: domainjogo.CreateJogoRequest{Nome: "   ", Tipo: "RPG", Nota: 9, Review: "Otimo"},
			wantErr: "nome e obrigatorio",
		},
		{
			name:    "tipo vazio",
			request: domainjogo.CreateJogoRequest{Nome: "Elden", Tipo: "", Nota: 9, Review: "Otimo"},
			wantErr: "tipo e obrigatorio",
		},
		{
			name:    "nota fora do intervalo",
			request: domainjogo.CreateJogoRequest{Nome: "Elden", Tipo: "RPG", Nota: 11, Review: "Otimo"},
			wantErr: "nota deve estar entre 1 e 10",
		},
		{
			name:    "review vazia",
			request: domainjogo.CreateJogoRequest{Nome: "Elden", Tipo: "RPG", Nota: 9, Review: ""},
			wantErr: "review e obrigatoria",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.request.Validate()
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}

			require.Error(t, err)
			assert.Equal(t, tc.wantErr, err.Error())
		})
	}
}

func TestCreateJogoRequestToEntityTrims(t *testing.T) {
	request := domainjogo.CreateJogoRequest{Nome: "  Elden  ", Tipo: " RPG ", Nota: 9, Review: " Otimo "}

	entity := request.ToEntity()

	assert.Equal(t, "Elden", entity.Nome)
	assert.Equal(t, "RPG", entity.Tipo)
	assert.Equal(t, 9, entity.Nota)
	assert.Equal(t, "Otimo", entity.Review)
}

func TestUpdateJogoRequestToEntityKeepsID(t *testing.T) {
	request := domainjogo.UpdateJogoRequest{Nome: "Elden", Tipo: "RPG", Nota: 9, Review: "Otimo"}

	entity := request.ToEntity(42)

	assert.Equal(t, 42, entity.ID)
}

func TestNewJogoResponseList(t *testing.T) {
	jogos := []domainjogo.Jogo{
		{ID: 1, Nome: "A"},
		{ID: 2, Nome: "B"},
	}

	response := domainjogo.NewJogoResponseList(jogos)

	require.Len(t, response, 2)
	assert.Equal(t, 1, response[0].ID)
	assert.Equal(t, "B", response[1].Nome)
}
