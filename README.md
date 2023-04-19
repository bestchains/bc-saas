<h1>
bc-saas represents software-as-a-service in bestchains which provides quick way to create a blockchain service
</h1>

Now we have:

- `digital depository` which provide basic capability to manage depositories

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
go build -o bin/depository cmd/depository/main.go
```

2. verify `depository`

```shell
❯ ./bin/depository -h
Usage of ./bin/depository:
  -addr string
    	used to listen and serve http requests (default ":9999")
  -profile string
    	profile to connect with blockchain network (default "./network.json")
  -contract string
    	contract name (default "depository")
  -v value
    	number for the log level verbosity
```

3. start `depository` server

```shell
./bin/depository -addr localhost:9999 -profile ./test/profile.json -contract depository
```

arguments:

- `-addr`: used to listen and serve http requests (default ":9999")
- `profile`: retrieved from [bestchains-platform](https://bestchains.github.io/website/docs/UserGuide/NetworkGov/channels)

## Development

### APIs

- `depository` APIS: [See the documentation](./doc/depository_api.md)

## Contribute to bc-saas

If you want to contribute to bc-saas,refer to [contribute guide](./CONTRIBUTING.md)

## Support

If you need support, start with the troubleshooting guide, or create github [issues](https://github.com/bestchains/bc-saas/issues/new)
