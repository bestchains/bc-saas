
## How to verify the success of listening to contract events

### Build a client
Create a simple command line tool with the following command.

```shll
git clone https://github.com/bestchains/bc-explorer.git

cd cmd/client

go build -o client main.go
```

### Start server

```shell
cd bc-saas/cmd/depository
go build main.go

./main -profile profile.json -contract depository -dsn 'postgres://bestchains:Passw0rd!@172.22.96.209:5432/bc-saas?sslmode=disable'
```

### Call the contract and confirming that the data is written to the database

1. call bc-saas PutUntrustValue
```shell
 curl -X POST \
  localhost:9999/basic/putUntrustValue \
  -H 'content-type: application/json' \
   -d '{
   "value":"eyJuYW1lIjoiYWJjIiwiY29udGVudFR5cGUiOiAianNvbiIsImNvbnRlbnRJRCI6ICJpZCIsInRydXN0ZWRUaW1lc3RhbXAiOiAiMTIzNCIsInBsYXRmb3JtIjogImJlc3RjaGFpbnMifQo="
}'
```

2. check db

```shell
bc-saas=> select * from "proof-c0zpw_depository_depository";
 index |                   kid                    |  platform  | operator | owner | blockNumber | contentName | contentID | contentType | trustedTimestamp 
-------+------------------------------------------+------------+----------+-------+-------------+-------------+-----------+-------------+------------------
 12    | a2be544f745022161bc1850053fc61e9ad4d0c8a | bestchains | a        | b     | 22          | abc         | id        | id          | 1234
(1 row)
```
