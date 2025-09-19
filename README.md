## E-commerce Stock Service (Go + Gin + Postgres)

Features:
- JWT auth (register/login)
- Product listing with available stock per shop
- Atomic checkout that reserves stock with row-level locks (no oversell)
- Idempotent POST /api/orders via Idempotency-Key
- Payment finalization that deducts inventory
- Warehouse activate/deactivate and transfer
- Background worker releasing expired reservations
- sqlx, validator, zap, middleware (request-id, logging, recovery, JWT)

### Quick start

1) Start Postgres and app
```bash
make docker-up
```

2) Run migrations (from host)
```bash
make migrate DB_URL="postgres://postgres:postgres@localhost:5432/ecommerce?sslmode=disable"
```

3) Health check
```bash
curl -s localhost:8080/api/healthz
```

### Auth
```bash
curl -s -X POST localhost:8080/api/register -H 'Content-Type: application/json' \
  -d '{"email":"a@b.com","password":"password123"}'

curl -s -X POST localhost:8080/api/login -H 'Content-Type: application/json' \
  -d '{"email":"a@b.com","password":"password123"}'
```

### Create order (idempotent)
```bash
curl -s -X POST localhost:8080/api/orders -H 'Content-Type: application/json' \
  -H 'Idempotency-Key: abc-123' \
  -d '{"shop_id":"<shop-uuid>","items":[{"product_id":"<product-uuid>","quantity":1}]}'
```

### Pay order
```bash
curl -s -X POST localhost:8080/api/orders/<order-id>/pay
```

### Warehouses
```bash
curl -s -X POST localhost:8080/api/warehouses/<id>/activate -H 'Authorization: Bearer <token>'
curl -s -X POST localhost:8080/api/warehouses/<id>/deactivate -H 'Authorization: Bearer <token>'
curl -s -X POST localhost:8080/api/warehouses/transfer -H 'Authorization: Bearer <token>' -H 'Content-Type: application/json' \
  -d '{"from":"<wh1>","to":"<wh2>","product_id":"<prod>","quantity":5}'
```

### Tests
```bash
make test
```



