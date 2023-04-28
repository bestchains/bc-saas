# Bestchains SaaS

[![codecov](https://codecov.io/gh/bestchains/bc-saas/branch/main/graph/badge.svg?token=PSLH95TQRS)](https://codecov.io/gh/bestchains/bc-saas)

bc-saas represents software-as-a-service provided in bestchains

Now we have:

- `digital depository` which provide basic capability to manage depositories

Service irrelevant architecture:

![saas_arch](./doc/images/arch.png)

## Usage

### Prerequsities

- [Golang](https://go.dev/)
- A blockchain network from [bestchains platform](https://bestchains.github.io/website/docs/QuickStart/usage)
- Deployed [depository contract](https://github.com/bestchains/bestchains-contracts/tree/main/examples/basic) on that network

### Build image

```bash
# output: hyperledgerk8s/bc-saas:v0.1.0
WHAT=bc-saas GOOS=linux GOARCH=amd64 make image
```

### Quick start

#### start a depository service

1. build `depository` server

```shell
WHAT=depository make binary
```

2. verify `depository`

```shell
❯ ./_output/bin/linux/amd64/depository -h
Usage of ./bin/depository:
  -auth string
        user authentication method, none, oidc or kubernetes (default "none")
  -addr string
     used to listen and serve http requests (default ":9999")
  -profile string
     profile to connect with blockchain network (default "./network.json")
  -contract string
     contract name (default "depository")
  -db string
        which database to use, default is pg(postgresql) (default "pg")
  -dsn string
        database connection string (default "postgres://bestchains:Passw0rd!@127.0.0.1:5432/bc-saas?sslmode=disable")
  -v value
     number for the log level verbosity
```

3. start `depository` server

```shell
./bin/depository -auth none -addr localhost:9999 -profile test/profile.json -contract depository -db pg -dsn 'postgres://bestchains:Passw0rd!@172.22.96.209:5432/bc-saas?sslmode=disable'
```

arguments:

- `-auth`: authentication method(default "none")
- `-addr`: used to listen and serve http requests (default ":9999")
- `-profile`: retrieved from [bestchains-platform](https://bestchains.github.io/website/docs/UserGuide/NetworkGov/channels)
- `-contract`: [depository contract's](https://github.com/bestchains/bestchains-contracts/tree/main/examples/depository) name deployed on the blockchain network
- `-db`: database type,only `pg` allowed for now
- `-dsn`: database connection string,only `postgresql` supported for now

> Only when `db` is `pg`,we will listen on contract events and store depository details into database

## Development

### APIs

- `depository` APIS: [See the documentation](./doc/depository_api.md)

## Contribute to bc-saas

If you want to contribute to bc-saas,refer to [contribute guide](./CONTRIBUTING.md)

## Support

If you need support, start with the troubleshooting guide, or create github [issues](https://github.com/bestchains/bc-saas/issues/new)
