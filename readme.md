# Inventory System

A full-stack inventory management system built with **Go** (backend) and **React** (frontend), backed by **PostgreSQL**. Built as a learning project covering REST API design, database transactions, optimistic locking, audit logging, and modern React data-fetching patterns.

---

## Features

- Add products with SKU, name, description, quantity, and price
- Stock in / out / adjustment operations with atomic transactions
- Optimistic locking (`SELECT FOR UPDATE`) to prevent race conditions on concurrent stock updates
- Low stock alerts with configurable per-product thresholds
- Inventory report — total value, total products, and low stock items
- Full-text product search and low-stock filter
- Audit log — every create, update, and delete records before/after values
- FK-protected deletes — products with transaction history cannot be removed
- JWT-based authentication middleware

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go, chi router |
| Database | PostgreSQL 16 |
| DB driver | pgx/v5 (connection pool) |
| Validation | go-playground/validator |
| Auth | golang-jwt/jwt |
| Frontend | React + Vite |
| Data fetching | TanStack Query |
| HTTP client | Axios |
| Styling | Tailwind CSS |
| Containerisation | Docker + Docker Compose |

---

## Project Structure

```
inventory/
├── backend/
│   ├── cmd/server/main.go          # Entry point, router setup
│   ├── internal/
│   │   ├── db/                     # Connection pool
│   │   ├── handlers/               # HTTP handlers (thin layer)
│   │   ├── middleware/             # Auth, transaction, validation
│   │   ├── models/                 # Structs: Product, StockTransaction, AuditLog
│   │   └── service/                # Business logic
│   ├── migrations/                 # SQL migration files
│   ├── Dockerfile
│   └── go.mod
└── frontend/
    ├── src/
    │   ├── api/                    # Axios client + TanStack Query hooks
    │   ├── components/             # ProductCard, StockModal, CreateProductForm, LoginGate
    │   └── pages/                  # ProductsPage, ReportPage, AuditPage
    └── package.json
```

---

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Node.js 18+ (for local frontend dev)
- Go 1.26+ (for local backend dev, optional if using Docker)
- Python 3 + `pyjwt` (for minting test tokens: `pip install pyjwt`)

### 1. Clone the repo

```bash
git clone https://github.com/karthikbhandary2/inventory-system.git
cd inventory-system
```

### 2. Start the containers

```bash
docker compose up -d
```

> If port 5432 is already in use by a local Postgres instance, change the host port mapping in `docker-compose.yml`:
> ```yaml
> postgres:
>   ports:
>     - "25432:5432"   # use any free host port
> ```
> The backend's `DB_PORT` env var should always remain `5432` — it talks to Postgres over the Docker network, not through the host mapping.

### 3. Run migrations

```bash
docker compose exec -T postgres psql -U inventory -d inventory < backend/migrations/001_create_products.sql
docker compose exec -T postgres psql -U inventory -d inventory < backend/migrations/002_stock_transactions.sql
docker compose exec -T postgres psql -U inventory -d inventory < backend/migrations/003_audit_log.sql
```

Verify:

```bash
docker compose exec postgres psql -U inventory -d inventory -c "\dt"
```

### 4. Mint a test JWT

There is no login endpoint yet. Generate a signed token manually using the same secret as `JWT_SECRET` in `docker-compose.yml`:

```bash
python3 -c "
import jwt, datetime
token = jwt.encode(
    {
        'user_id': 'test-user-id',
        'username': 'tester',
        'exp': datetime.datetime.utcnow() + datetime.timedelta(hours=24)
    },
    'dev-secret-change-in-prod',
    algorithm='HS256'
)
print(token)
"
```

Export it for curl testing:

```bash
export TOKEN="<paste token here>"
```

### 5. Start the frontend

```bash
cd frontend
npm install
npm run dev
```

Visit `http://localhost:5173`, paste your JWT into the login gate, and the app connects to `http://localhost:8080/api/v1`.

---

## API Reference

All routes require `Authorization: Bearer <token>`.

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/products` | List products. Query params: `search`, `low_stock=true` |
| `POST` | `/api/v1/products` | Create a product |
| `GET` | `/api/v1/products/:id` | Get a single product |
| `PUT` | `/api/v1/products/:id` | Update name, description, price, threshold |
| `DELETE` | `/api/v1/products/:id` | Delete (fails with 409 if stock transactions exist) |
| `POST` | `/api/v1/products/:id/stock` | Stock in / out / adjustment |
| `GET` | `/api/v1/reports/inventory` | Total value, total products, low stock items |
| `GET` | `/api/v1/audit` | Audit log. Query param: `entity_id` to filter |

### Stock operation body

```json
{
  "operation": "in",
  "quantity": 20,
  "notes": "weekly restock",
  "performed_by": "tester"
}
```

`operation` must be one of `in`, `out`, or `adjustment`. An `out` operation that would result in negative stock returns `409 Conflict`.

---

## Key Concepts Implemented

**Transaction middleware** — every route under `/api/v1` runs inside a PostgreSQL transaction. If the handler returns a 4xx/5xx status, the transaction rolls back automatically. This means a failed stock withdrawal never partially updates the quantity.

**Optimistic locking** — `StockOperation` uses `SELECT ... FOR UPDATE` to lock the product row before reading its quantity, preventing two concurrent requests from both reading the same value and double-deducting stock.

**Sentinel errors** — `ErrInsufficientStock` is a typed sentinel wrapped with `%w`, allowing handlers to use `errors.Is()` for precise status code mapping (409 vs 500) without fragile string matching.

**Audit log** — written inside the same transaction as the change, so the log entry and the change always succeed or fail together. Stores full JSON snapshots of before/after state.

**Database constraints as a safety net** — `CHECK (quantity >= 0)` on the products table ensures stock never goes negative even if application-level validation has a bug.

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DB_USER` | `inventory` | Postgres username |
| `DB_PASS` | `secret` | Postgres password |
| `DB_HOST` | `postgres` | Postgres hostname (Docker service name) |
| `DB_PORT` | `5432` | Postgres port (internal container port) |
| `DB_NAME` | `inventory` | Database name |
| `JWT_SECRET` | `dev-secret-change-in-prod` | HMAC secret for JWT signing — **change in production** |

---

## Rebuilding after code changes

Docker caches build layers. After editing Go source files, force a full rebuild:

```bash
docker compose build backend --no-cache
docker compose up -d
```

`docker compose restart` does **not** rebuild the image — always use `build` + `up` after source changes.

---

## Running API tests manually

```bash
# Create a product
curl -i -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"sku":"WIDGET-001","name":"Blue Widget","quantity":50,"price":9.99,"low_stock_threshold":10}'

export PRODUCT_ID="<id from response>"

# Stock in
curl -i -X POST http://localhost:8080/api/v1/products/$PRODUCT_ID/stock \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"operation":"in","quantity":20,"notes":"restock","performed_by":"tester"}'

# Over-withdrawal — expect 409
curl -i -X POST http://localhost:8080/api/v1/products/$PRODUCT_ID/stock \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"operation":"out","quantity":9999,"notes":"too much","performed_by":"tester"}'

# Report
curl -s http://localhost:8080/api/v1/reports/inventory -H "Authorization: Bearer $TOKEN" | jq

# Audit log
curl -s "http://localhost:8080/api/v1/audit?entity_id=$PRODUCT_ID" -H "Authorization: Bearer $TOKEN" | jq
```

---

## What's next

- `/auth/login` endpoint to replace the manual JWT paste flow
- Role-based access control (admin vs read-only)
- Pagination on the products list and audit log
- Stock transaction history per product
- CSV export for the inventory report