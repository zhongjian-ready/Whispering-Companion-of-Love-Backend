# WeChat MiniApp Backend

This is a backend service for a WeChat Mini Program, written in Go.

## Tech Stack

- **Language**: Go
- **Web Framework**: Gin
- **Database**: MySQL (using GORM)
- **Cache**: Redis (using go-redis)
- **Configuration**: Viper

## Project Structure

- `main.go`: Entry point of the application.
- `internal/config`: Configuration loading logic.
- `internal/router`: HTTP router setup.
- `pkg/database`: Database connection setup.
- `pkg/cache`: Redis connection setup.

## Configuration

Configuration is managed via `config.yaml`. You can also use environment variables.

## Running Locally

1. Ensure MySQL and Redis are running.
2. Update `config.yaml` with your database and redis credentials.
3. Run the server:
   ```bash
   go run main.go
   ```

## Dockerfile 最佳实践

请参考[如何提高项目构建效率](https://developers.weixin.qq.com/miniprogram/dev/wxcloudrun/src/scene/build/speed.html)
