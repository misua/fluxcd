# Go CRUD API

A simple REST API built with Go that demonstrates CRUD operations for FluxCD deployment.

## Features

- **In-memory storage** (no database required)
- **RESTful endpoints** for items management
- **Health check** endpoint
- **Environment-aware** (dev/staging/production)
- **Lightweight** Docker image (~15MB)

## API Endpoints

### Health Check
```bash
GET /health
```

Response:
```json
{
  "status": "healthy",
  "environment": "development",
  "timestamp": "2025-10-28T14:00:00Z"
}
```

### Get All Items
```bash
GET /items
```

### Get Item by ID
```bash
GET /items/{id}
```

### Create Item
```bash
POST /items
Content-Type: application/json

{
  "name": "Example Item",
  "description": "This is an example"
}
```

### Update Item
```bash
PUT /items/{id}
Content-Type: application/json

{
  "name": "Updated Item",
  "description": "Updated description"
}
```

### Delete Item
```bash
DELETE /items/{id}
```

## Local Development

```bash
# Install dependencies
go mod download

# Run locally
go run main.go

# Test the API
curl http://localhost:8080/health
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Item","description":"A test item"}'
curl http://localhost:8080/items
```

## Docker Build

```bash
# Build image
docker build -t go-crud-api:latest .

# Run container
docker run -p 8080:8080 -e ENVIRONMENT=development go-crud-api:latest

# Test
curl http://localhost:8080/health
```

## Environment Variables

- `PORT` - Server port (default: 8080)
- `ENVIRONMENT` - Environment name (dev/staging/production)

## Kubernetes Deployment

See the `fluxcd-apps` repository for Kubernetes manifests and FluxCD configuration.
