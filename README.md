# Oauth starter

An [application continuum](https://www.appcontinuum.io/) style example using golang
that includes an oauth client, authorization server, and resource server.

## Getting Started

_hang tight!_

This project is work in progress.

## Development

Generate private and public keys.

```bash
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -pubout -outform PEM -out public.pem
```

That's a wrap for now.