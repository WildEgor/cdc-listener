# Base Stage
FROM golang:1.22-alpine AS base
LABEL maintainer="YOUR_NAME <YOUR_EMAIL>"
# if use private libs uncomment this
#ARG GITHUB_TOKEN
#RUN apk update && apk add ca-certificates git openssh
#RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && mkdir -p dist

# Development Stage
FROM base as dev
RUN go install github.com/cosmtrek/air@latest
WORKDIR /app/
COPY . .
RUN go build -o dist/app cmd/main.go
CMD ["air", "-c", ".air-unix.toml", "-d"]

# Build Production Stage
FROM base as builder
WORKDIR /app
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN $GOPATH/bin/swag init -g cmd/main.go --output docs
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o dist/app cmd/main.go

# Production Stage
FROM cgr.dev/chainguard/busybox:latest-glibc as production
WORKDIR /app/
COPY --from=builder /app/docs/* ./docs
COPY --from=builder /app/dist/app .
# Specify method fetch .env!
COPY --from=builder /app/.env.local .
RUN mkdir "data"
CMD ["./app"]