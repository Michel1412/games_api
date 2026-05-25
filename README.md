# games_api

API REST em Go para um CRUD simples de reviews de jogos, pensada para servir como backend de um aplicativo mobile.

## Stack

- Go 1.22+ (o `go.mod` declara `go 1.25.0`; basta ter uma toolchain compativel instalada)
- `net/http` + [`chi`](https://github.com/go-chi/chi)
- CORS via [`go-chi/cors`](https://github.com/go-chi/cors)
- Firestore em producao
- Repositorio in-memory em desenvolvimento/testes
- Testes com `testing` e [`testify`](https://github.com/stretchr/testify)

## Estrutura

```
.
├── cmd/server/             # entrypoint da aplicacao (package main)
├── internal/
│   ├── bootstrap/          # construcao do repositorio (memoria/Firestore)
│   ├── config/             # leitura de variaveis de ambiente
│   ├── domain/             # entidades e DTOs (jogo, user)
│   ├── handler/            # rotas HTTP, CORS, handlers
│   ├── service/            # implementacoes de JogoRepository
│   └── usecase/            # regras de negocio (jogo, user)
└── tests/                  # toda a suite de testes, espelhando a estrutura interna
    ├── bootstrap/
    ├── config/
    ├── domain/{jogo,user}/
    ├── handler/
    ├── service/
    └── usecase/{jogo,user}/
```

Os testes ficam **isolados em `tests/`** em vez de espalhados como `*_test.go` ao lado do codigo de producao. Cada pacote de teste usa o sufixo `_test` (caixa preta) e importa apenas a API publica de `internal/...`.

## Rodando localmente

```bash
go run ./cmd/server
```

Por padrao a API sobe em `http://localhost:8080` usando `ENV=dev` e o repositorio in-memory. Para construir um binario:

```bash
go build -o bin/server ./cmd/server
./bin/server
```

## Variaveis de ambiente

Copie `.env.example` para `.env` (apenas referencia, o app le do ambiente) e defina:

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

## CORS

CORS configurado em `internal/handler/router.go` permite qualquer origem, qualquer header e os metodos `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`. Como a API autentica via token Bearer (e nao por cookies), `AllowCredentials` esta desabilitado, mantendo compatibilidade com `Access-Control-Allow-Origin: *`.

## Testes

A suite vive em `tests/` e pode ser executada com:

```bash
# Apenas executar a suite
go test ./tests/...

# Com saida detalhada
go test ./tests/... -v

# Com cobertura agregada sobre o codigo de producao
go test ./tests/... -coverpkg=./internal/... -coverprofile=coverage.out
go tool cover -func=coverage.out          # resumo no terminal
go tool cover -html=coverage.out          # relatorio HTML no navegador
```

> Note que a cobertura usa `-coverpkg=./internal/...` porque os testes estao em arvore separada. Sem essa flag a cobertura sairia como 0% (os pacotes de teste nao contem codigo de producao). O pacote `cmd/server` ficou propositadamente fino — toda a logica de bootstrap esta em `internal/bootstrap/` e e testada de la.

Rodar um pacote especifico:

```bash
go test ./tests/handler -v
go test ./tests/usecase/jogo -run TestCreateJogoUseCaseCreatesJogo -v
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
