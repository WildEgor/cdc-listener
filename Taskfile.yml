version: '3'

tasks:
  # Run using docker
  docker-dev:
    dir: '{{.USER_WORKING_DIR}}'
    preconditions:
      - test -f docker-compose.yml
    cmd: docker-compose up -d app-dev

  local-dev:
    cmd: air -d

  # Update DI deps
  wire:
    cmd: cd internal && wire

  # Update swagger files
  swag:
    cmd: swag init -g cmd/main.go --output docs