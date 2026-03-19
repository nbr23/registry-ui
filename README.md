# registry-ui

A minimal web UI for docker [`registry`](https://hub.docker.com/_/registry).

## Running with Docker Compose

**1. Create an htpasswd file:**

```sh
htpasswd -Bc .htpasswd yourusername
```

**2. Start the stack:**

```sh
docker compose up -d
```

The UI will be available at `http://localhost:8080`. Log in with the credentials from your `.htpasswd` file.

## Pushing an image

```sh
docker login localhost:5000
docker pull alpine
docker tag alpine localhost:5000/alpine:latest
docker push localhost:5000/alpine:latest
```

## Configuration

| Variable       | Default                   | Description                        |
|----------------|---------------------------|------------------------------------|
| `REGISTRY_URL` | `http://localhost:5000`   | URL of the backend Docker registry |
| `PULL_HOST`    | host from `REGISTRY_URL`  | Hostname shown in pull commands (override if your registry is behind a reverse proxy) |

## Building manually

```sh
go build -o registry-ui .
REGISTRY_URL=http://localhost:5000 ./registry-ui
```
