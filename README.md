# Clean Architecture - Orders

Aplicação Go com Clean Architecture expondo a listagem de orders via REST, gRPC e GraphQL.

## Execução

```bash
docker compose up
```

## Portas

- REST: http://localhost:8080/order
- gRPC: localhost:50051
- GraphQL: http://localhost:8081/graphql

## Exemplo de query GraphQL

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
