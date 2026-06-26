# Processador de Faturas Assíncrono

API desenvolvida utilizando Go, MongoDB e RabbitMQ que simula o processamento de faturas de forma assíncrona com workers orientados a mensagens.

## Tecnologias

- **Go 1.24.5** -- linguagem principal
- **Gin** -- framework HTTP
- **MongoDB 6.0** -- banco de dados (colecao `faturas`)
- **RabbitMQ** -- message broker (filas: `save`, `generateNote`, `notifications`)

## Requisitos

- Docker e Docker Compose
- Arquivo `.env` na raiz do projeto (veja `.env.example`)

## Executando

```bash
docker compose up --build
```

A API fica disponivel em `http://localhost:1234`.

Servicos auxiliares:
- Mongo Express (interface MongoDB): `http://localhost:8081`
- RabbitMQ Management UI: `http://localhost:15672`

## Endpoints

### POST /faturas

Envia uma lista de faturas para processamento assincrono. Cada fatura e validada e enfileirada no RabbitMQ para salvamento, emissao de nota fiscal e notificacao.

**Corpo da requisicao:** array de objetos `Fatura`

```json
[
  {
    "id": "FAT-20250821-101",
    "cnpj": "69.375.897/0001-34",
    "valorTotal": 10,
    "descricao": "Desenvolvimento de API"
  }
]
```

**Validacoes por item:**
- `cnpj` -- deve ser um CNPJ valido (utiliza a biblioteca `brdoc`)
- `valorTotal` -- deve ser maior que zero
- `descricao` -- nao pode ser vazia (espacos sao desconsiderados)

**Resposta de sucesso (200):**
```json
{
  "faturas": [
    {
      "id": "FAT-20250821-101",
      "cnpj": "69.375.897/0001-34",
      "valorTotal": 10,
      "descricao": "Desenvolvimento de API"
    }
  ]
}
```

**Resposta de erro (500):** retorna um objeto com o campo `error` descrevendo o problema. A requisicao e rejeitada por completo se algum item falhar na validacao ou no enfileiramento.

---

### GET /faturas

Retorna todas as faturas salvas no MongoDB.

**Resposta de sucesso (200):**
```json
{
  "faturas": [
    {
      "id": "FAT-20250821-101",
      "cnpj": "69.375.897/0001-34",
      "valorTotal": 10,
      "descricao": "Desenvolvimento de API"
    }
  ]
}
```

**Resposta de erro (500):** retorna um objeto com o campo `error` em caso de falha na consulta ao banco.

## Estrutura do projeto

```
main.go                     entrada da aplicacao
internal/
  database/                 conexao com MongoDB (banco: faturasAPI, colecao: faturas)
  dates/                    utilitarios de formatacao de data
  fatura/                   logica de negocios, handlers HTTP, workers, repositorio MongoDB
  logger/                   escrita de logs em arquivos (tmp/<YYYY-M-D>.txt)
  message/                  conexao e operacoes com RabbitMQ
```

O fluxo de processamento segue a sequencia de filas no RabbitMQ:

```
POST /faturas -> validacao -> [fila save] -> [fila generateNote] -> [fila notifications] -> log em tmp/
```

Quatro workers rodam em segundo plano (goroutines): dois `SaveWorker`, um `GenerateNoteWorker` e um `NotificationsWorker`.
