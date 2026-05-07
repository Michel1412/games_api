# games_api

API REST em Go para um CRUD simples de reviews de jogos, pensada para servir como backend de um aplicativo mobile.

## Stack

- Go 1.25
- `net/http` + `chi`
- Firestore em producao
- Repositorio in-memory em desenvolvimento/testes
- Testes com `testing` e `testify`

## Rodando localmente

```bash
go run ./cmd/server
```

Por padrao a API sobe em `http://localhost:8080` usando `ENV=dev` e o repositorio in-memory.

## Variaveis de ambiente

```bash
PORT=8080
ENV=dev
GCP_PROJECT_ID=
GOOGLE_APPLICATION_CREDENTIALS=
GOOGLE_APPLICATION_CREDENTIALS_JSON=
```

Use `ENV=prod` para ativar Firestore. Em producao, informe `GCP_PROJECT_ID` e uma das credenciais:

- `GOOGLE_APPLICATION_CREDENTIALS`: caminho para o arquivo JSON da service account.
- `GOOGLE_APPLICATION_CREDENTIALS_JSON`: conteudo completo do JSON, util para Railway.

## Endpoints

### POST `/login`

Credenciais validas:

```json
{
  "email": "usuario@esoft.com",
  "password": "Abc123"
}
```

Resposta `200 OK`:

```json
{
  "token": "550e8400-e29b-41d4-a716-446655440000"
}
```

### GET `/jogos`

Retorna todos os jogos.

### GET `/jogos/{id}`

Retorna um jogo pelo ID sequencial.

### POST `/jogos`

```json
{
  "nome": "Elden Ring",
  "tipo": "RPG",
  "nota": 9,
  "review": "Desafiador e visualmente impecavel."
}
```

### PUT `/jogos/{id}`

Atualiza todos os dados do jogo. Todos os campos sao obrigatorios.

### DELETE `/jogos/{id}`

Remove fisicamente o jogo.

## Validacoes e erros

- `nome`, `tipo` e `review` sao obrigatorios.
- `nota` deve estar entre 1 e 10.
- Erros retornam `400 Bad Request`:

```json
{
  "error": "nota deve estar entre 1 e 10"
}
```

## Testes

```bash
go test ./... -cover
```

## Deploy no Railway

O deploy deve ser feito somente pela branch `main`:

1. Crie um novo servico no Railway apontando para este repositorio.
2. Em `Service > Settings > Source`, defina `Production Branch` como `main`.
3. Desative deploys de PR/preview se nao quiser ambientes temporarios.
4. Configure as variaveis:
   - `ENV=prod`
   - `PORT` deixado pelo Railway ou `8080`
   - `GCP_PROJECT_ID`
   - `GOOGLE_APPLICATION_CREDENTIALS_JSON` com o JSON da service account
5. O Railway usara o `Dockerfile` via `railway.json`.
