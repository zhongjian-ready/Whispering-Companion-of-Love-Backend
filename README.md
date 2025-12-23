# WeChat MiniApp Backend

This is a backend service for a WeChat Mini Program, written in Go.

## Tech Stack

- **Language**: Go
- **Web Framework**: Gin
- **Database**: MySQL (using GORM)
- **Configuration**: Viper

## Project Structure

- `cmd/server/main.go`: Entry point of the application.
- `internal/config`: Configuration loading logic.
- `internal/router`: HTTP router setup.
- `internal/handler`: HTTP request handlers.
- `internal/model`: Database models.
- `internal/repository`: Database access layer.
- `pkg/database`: Database connection setup.

## Configuration

Configuration is managed via `config.yaml`. You can also use environment variables.

## Running Locally

1. Ensure MySQL is running.
2. Update `config.yaml` with your database credentials.
3. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

## Dockerfile 最佳实践

请参考[如何提高项目构建效率](https://developers.weixin.qq.com/miniprogram/dev/wxcloudrun/src/scene/build/speed.html)
