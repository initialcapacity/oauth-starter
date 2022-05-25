# OAuth starter

An [application continuum](https://www.appcontinuum.io/) style example using golang
that includes an oauth client, authorization server, and resource server.

## Getting Started

Install the following prerequisites.

* [Go 1.18](https://go.dev)
* [Pack](https://buildpacks.io)
* [Docker Desktop](https://www.docker.com/products/docker-desktop)

Build with Pack.

```bash
pack build oauth-starter --builder heroku/buildpacks:20
```

Run with docker compose.

```bash
docker-compose up
````

## Development

Generate private and public keys.

```bash
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -pubout -outform PEM -out public.pem
```

That's a wrap for now.