# WeChat MiniApp Backend

This is a backend service for a WeChat Mini Program, written in Go.

## Tech Stack

- **Language**: Go
- **Web Framework**: Gin
- **Database**: PostgreSQL (using GORM)
- **Cache**: Redis (using go-redis)
- **Configuration**: Viper

## Project Structure

- `cmd/server`: Entry point of the application.
- `internal/config`: Configuration loading logic.
- `internal/router`: HTTP router setup.
- `pkg/database`: Database connection setup.
- `pkg/cache`: Redis connection setup.

## Configuration

Configuration is managed via `config.yaml`. You can also use environment variables.

## Running Locally

1. Ensure PostgreSQL and Redis are running.
2. Update `config.yaml` with your database and redis credentials.
3. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

## Deployment (Docker)

Build the Docker image:

```bash
docker build -t miniapp-backend .
```

Run the container:

```bash
docker run -p 8080:8080 miniapp-backend
```
