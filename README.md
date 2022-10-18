# Homebugh API


An API for the [Homebugh](https://github.com/ck3g/homebugh) application.

## Configure development environment

1. Generate a self-signed TLS certificate

```shell
$ cd tls
$ go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

when using asdf with 1.19.2 version of go (adjust if needed):
```
go run ~/.asdf/installs/golang/1.19.2/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

2. Configure database

The application uses the database from [Homebugh](https://github.com/ck3g/homebugh) repository.
Follow the steps from [Configure development environment](https://github.com/ck3g/homebugh#configure-development-environment) list.

## License

HomeBugh app is released under the [MIT License](./LICENSE).
