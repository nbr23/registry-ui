# registry-ui

A minimal web UI for docker [`registry`](https://hub.docker.com/_/registry) Docker registry.

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

## Configuration

| Variable       | Default                   | Description                        |
|----------------|---------------------------|------------------------------------|
| `REGISTRY_URL` | `http://localhost:5000`   | URL of the backend Docker registry |

## Building manually

```sh
go build -o registry-ui .
REGISTRY_URL=http://localhost:5000 ./registry-ui
```
