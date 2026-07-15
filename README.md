# Clean Architecture - Orders

Aplicação Go com Clean Architecture que expõe operações de pedidos (`orders`) via REST, gRPC e GraphQL.

## Visão geral

O projeto contém:
- `REST` em `http://localhost:8080/order`
- `gRPC` em `localhost:50051`
- `GraphQL` em `http://localhost:8081/graphql`
- Banco de dados PostgreSQL gerenciado via `docker compose`
- Migrações de banco de dados automatizadas na inicialização

## Pré-requisitos

- Docker
- Docker Compose
- Go 1.22+ (opcional, para rodar localmente sem Docker)

## Executando a aplicação

### 1) Usando Docker Compose

No diretório raiz do projeto, execute:

```bash
docker compose up --build
```

Isso iniciará os serviços:
- `db`: PostgreSQL com banco `orders`
- `app`: aplicação Go

### 2) Parar os serviços

```bash
docker compose down
```

## Variáveis de ambiente

O serviço `app` usa as seguintes variáveis, definidas em `docker-compose.yaml`:

- `DB_HOST`: host do PostgreSQL (`db` no Docker Compose)
- `DB_PORT`: porta do PostgreSQL (`5432`)
- `DB_USER`: usuário do banco (`orders`)
- `DB_PASSWORD`: senha do banco (`orders`)
- `DB_NAME`: nome do banco (`orders`)
- `DB_SSLMODE`: SSL mode do Postgres (`disable`)
- `REST_ADDR`: endereço do servidor REST (`:8080`)
- `GRPC_ADDR`: endereço do servidor gRPC (`:50051`)
- `GRAPHQL_ADDR`: endereço do servidor GraphQL (`:8081`)
- `MIGRATIONS_PATH`: caminho das migrações dentro do container

## Endpoints disponíveis

### REST

- Listar orders: `GET http://localhost:8080/order`
- Criar order: `POST http://localhost:8080/order`

Exemplo de requisição GET:

```bash
curl http://localhost:8080/order
```

### gRPC

- Endereço: `localhost:50051`
- Serviço: `OrderService`
- Método: `ListOrders`

Exemplo usando `grpcurl`:

```bash
grpcurl -plaintext localhost:50051 orderspb.OrderService/ListOrders
```

### GraphQL

- Endpoint: `http://localhost:8081/graphql`

Exemplo de query:

```graphql
query {
  listOrders {
    id
    customerName
    total
    status
  }
}
```

## Testando a aplicação

1. Execute `docker compose up --build`.
2. Aguarde até o serviço `app` iniciar e o PostgreSQL ficar disponível.
3. Acesse os endpoints REST, gRPC ou GraphQL.
4. Para criar pedidos via REST, envie um `POST` com JSON ao endpoint `/order`.

Exemplo de criação via `curl`:

```bash
curl -X POST http://localhost:8080/order \
  -H "Content-Type: application/json" \
  -d '{"customerName":"João","total":100.0,"status":"PENDING"}'
```

## Considerações finais

- As migrações são executadas automaticamente na inicialização da aplicação.
- O projeto está organizado segundo os princípios de Clean Architecture, com camadas de domínio, repositório, transportes e casos de uso.
- Se desejar rodar sem Docker, configure o PostgreSQL localmente e exporte as mesmas variáveis de ambiente usadas pelo `docker-compose.yaml`.
