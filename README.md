# Homebugh API


An API for the [Homebugh](https://github.com/ck3g/homebugh) application.

## Configure development environment

1. Generate a self-signed TLS certificate

```shell
$ cd tls
$ go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```