# OAuth starter

[![Build results](https://github.com/initialcapacity/oauth-starter/workflows/build/badge.svg)](https://github.com/initialcapacity/oauth-starter/actions)
[![codecov](https://codecov.io/gh/initialcapacity/oauth-starter/branch/main/graph/badge.svg)](https://codecov.io/gh/initialcapacity/oauth-starter)

An [application continuum](https://www.appcontinuum.io/) style example using Golang
that includes an OAuth 2 server with PKCE support.

* [OAuth 2.0 Authorization Framework](https://datatracker.ietf.org/doc/html/rfc6749)
* [Proof Key for Code Exchange](https://datatracker.ietf.org/doc/html/rfc7636)
* [OpenID Connect](https://openid.net/specs/openid-connect-core-1_0.html)

## Getting Started

Install the following prerequisites.

* [Go 1.18](https://go.dev)
* [Pack](https://buildpacks.io)
* [Docker Desktop](https://www.docker.com/products/docker-desktop)
* [Node.js](https://nodejs.org/en/)

Build golang binaries with Pack.

```bash
pack build oauth-starter --builder heroku/buildpacks:20
```

Build and pack the frontend app with NPM.

```bash
cd web
npm run build
npm run pack
cd -
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