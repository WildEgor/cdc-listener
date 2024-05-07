# Project boilerplate: Go + Fiber + Docker + Air
A containerised eShopFiberMicroserviceBoilerplate (mmm Java-style naming)

| Component         | Link                                   |
| ----------------- |----------------------------------------|
| Language          | [Go](https://go.dev/)                  |
| Framework         | [Fiber](https://docs.gofiber.io/)      |
| Hot Reloading     | [Air](https://github.com/cosmtrek/air) |
| Containerisation  | [Docker](https://www.docker.com/)      |

## Structure

- [/api](api/README.md)
- [/assets](assets/README.md)
- [/cmd](cmd/README.md)
- [/internal/adapters](internal/adapters/README.md)
- [/internal/configs](internal/configs/README.md)
- [/internal/handlers](internal/handlers/README.md)
- [/internal/mappers](internal/mappers/README.md)
- [/internal/middlewares](internal/middlewares/README.md)
- [/internal/models](internal/models/README.md)
- [/internal/repositories](internal/repositories/README.md)
- [/internal/routers](internal/routers/README.md)

## Development

Install the dependencies:

> This project uses [Go mod](https://blog.golang.org/using-go-modules), the official module manager, to handle Go modules in a portable way without having to worry about GOPATH.

```bash
go mod download
go mod vendor
go mod verify
```

Define environment variables for your development environment:

> These are passed to the Docker container via `docker-compose.yaml` in development. When running in production, the environment variables must be passed to the container when it is run.

```bash
cp .env.example .env
```

Run locally:

> This builds the Docker image and runs it automatically with the config defined in `docker-compose.yaml`. This saves you having to build the docker image and then run a manual `docker run` command with all the flags (for environment variables, ports, etc).

Local
```bash
task wire && air
```

Using Docker
```bash
docker compose up --build --remove-orphans app-dev
```

## Production

> Note: Environment variables are never baked into the image, or they wouldn't be _environment_ variables. The production environment will start a Docker container based on this image, but it will have to pass the environment variables to the container when it runs it.

Example manually running a container with environment variables and ports defined:
Using Docker
```bash
docker compose up --build --remove-orphans app
```