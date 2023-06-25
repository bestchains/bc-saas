# Depository APIs

### GET /basic/nonce

Used to get current nonce of a account

#### Example

```shell
curl -X GET \
  http://localhost:9999/basic/currentNonce?account=xxx 
```

#### Response

```json
{
  "nonce": "xxx"
}
```

### GET /basic/total/hf/metadata

Used to get depository contract's metadta

#### Example

```shell
curl -X GET \
  http://localhost:9999/hf/metadata 
```

#### Response

```json
{
    "content": "eyJpbmZvIjp7InRpdGxlIjoidW5kZWZpbmVkIiwidmVyc2lvbiI6ImxhdGVzdCJ9LCJjb250cmFjdHMiOnsib3JnLmJlc3RjaGFpbnMuY29tLkJhc2ljQ29udHJhY3QiOnsiaW5mbyI6eyJ0aXRsZSI6Im9yZy5iZXN0Y2hhaW5zLmNvbS5CYXNpY0NvbnRyYWN0IiwidmVyc2lvbiI6ImxhdGVzdCJ9LCJuYW1lIjoib3JnLmJlc3RjaGFpbnMuY29tLkJhc2ljQ29udHJhY3QiLCJ0cmFuc2FjdGlvbnMiOlt7InBhcmFtZXRlcnMiOlt7Im5hbWUiOiJwYXJhbTAiLCJzY2hlbWEiOnsidHlwZSI6InN0cmluZyJ9fSx7Im5hbWUiOiJwYXJhbTEiLCJzY2hlbWEiOnsidHlwZSI6Im51bWJlciIsImZvcm1hdCI6ImRvdWJsZSIsIm1heGltdW0iOjE4NDQ2NzQ0MDczNzA5NTUyMDAwLCJtaW5pbXVtIjowLCJtdWx0aXBsZU9mIjoxfX1dLCJ0YWciOlsic3VibWl0IiwiU1VCTUlUIl0sIm5hbWUiOiJDaGVjayJ9LHsicGFyYW1ldGVycyI6W3sibmFtZSI6InBhcmFtMCIsInNjaGVtYSI6eyJ0eXBlIjoic3RyaW5nIn19XSwidGFnIjpbInN1Ym1pdCIsIlNVQk1JVCJdLCJuYW1lIjoiQ3VycmVudCIsInJldHVybnMiOnsidHlwZSI6Im51bWJlciIsImZvcm1hdCI6ImRvdWJsZSIsIm1heGltdW0iOjE4NDQ2NzQ0MDczNzA5NTUyMDAwLCJtaW5pbXVtIjowLCJtdWx0aXBsZU9mIjoxfX0seyJ0YWciOlsic3VibWl0IiwiU1VCTUlUIl0sIm5hbWUiOiJEaXNhYmxlQUNMIn0seyJ0YWciOlsic3VibWl0IiwiU1VCTUlUIl0sIm5hbWUiOiJFbmFibGVBQ0wifSx7InBhcmFtZXRlcnMiOlt7Im5hbWUiOiJwYXJhbTAiLCJzY2hlbWEiOnsidHlwZSI6ImFycmF5IiwiaXRlbXMiOnsidHlwZSI6ImludGVnZXIiLCJmb3JtYXQiOiJpbnQzMiIsIm1heGltdW0iOjI1NSwibWluaW11bSI6MH19fV0sInRhZyI6WyJzdWJtaXQiLCJTVUJNSVQiXSwibmFtZSI6IkdldFJvbGVBZG1pbiIsInJldHVybnMiOnsidHlwZSI6ImFycmF5IiwiaXRlbXMiOnsidHlwZSI6ImludGVnZXIiLCJmb3JtYXQiOiJpbnQzMiIsIm1heGltdW0iOjI1NSwibWluaW11bSI6MH19fSx7InBhcmFtZXRlcnMiOlt7Im5hbWUiOiJwYXJhbTAiLCJzY2hlbWEiOnsidHlwZSI6InN0cmluZyJ9fV0sInRhZyI6WyJzdWJtaXQiLCJTVUJNSVQiXSwibmFtZSI6IkdldFZhbHVlQnlJbmRleCIsInJldHVybnMiOnsidHlwZSI6InN0cmluZyJ9fSx7InBhcmFtZXRlcnMiOlt7Im5hbWUiOiJwYXJhbTAiLCJzY2hlbWEiOnsidHlwZSI6InN0cmluZyJ9fV0sInRhZyI6WyJzdWJtaXQiLCJTVUJNSVQiXSwibmFtZSI6IkdldFZhbHVlQnlLSUQiLCJyZXR1cm5zIjp7InR5cGUiOiJzdHJpbmcifX0seyJwYXJhbWV0ZXJzIjpbeyJuYW1lIjoicGFyYW0wIiwic2NoZW1hIjp7InR5cGUiOiJhcnJheSIsIml0ZW1zIjp7InR5cGUiOiJpbnRlZ2VyIiwiZm9ybWF0IjoiaW50MzIiLCJtYXhpbXVtIjoyNTUsIm1pbmltdW0iOjB9fX0seyJuYW1lIjoicGFyYW0xIiwic2NoZW1hIjp7InR5cGUiOiJzdHJpbmcifX1dLCJ0YWciOlsic3VibWl0IiwiU1VCTUlUIl0sIm5hbWUiOiJHcmFudFJvbGUifSx7InBhcmFtZXRlcnMiOlt7Im5hbWUiOiJwYXJhbTAiLCJzY2hlbWEiOnsidHlwZSI6ImFycmF5IiwiaXRlbXMiOnsidHlwZSI6ImludGVnZXIiLCJmb3JtYXQiOiJpbnQzMiIsIm1heGltdW0iOjI1NSwibWluaW11bSI6MH19fSx7Im5hbWUiOiJwYXJhbTEiLCJzY2hlbWEiOnsidHlwZSI6InN0cmluZyJ9fV0sInRhZyI6WyJzdWJtaXQiLCJTVUJNSVQiXSwibmFtZSI6Ikhhc1JvbGUiLCJyZXR1cm5zIjp7InR5cGUiOiJib29sZWFuIn19LHsicGFyYW1ldGVycyI6W3sibmFtZSI6InBhcmFtMCIsInNjaGVtYSI6eyJ0eXBlIjoic3RyaW5nIn19XSwidGFnIjpbInN1Ym1pdCIsIlNVQk1JVCJdLCJuYW1lIjoiSW5jcmVtZW50IiwicmV0dXJucyI6eyJ0eXBlIjoibnVtYmVyIiwiZm9ybWF0IjoiZG91YmxlIiwibWF4aW11bSI6MTg0NDY3NDQwNzM3MDk1NTIwMDAsIm1pbmltdW0iOjAsIm11bHRpcGxlT2YiOjF9fSx7InRhZyI6WyJzdWJtaXQiLCJTVUJNSVQiXSwibmFtZSI6IkluaXRpYWxpemUifSx7InRhZyI6WyJzdWJtaXQiLCJTVUJNSVQiXSwibmFtZSI6Ik93bmVyIiwicmV0dXJucyI6eyJ0eXBlIjoic3RyaW5nIn19LHsicGFyYW1ldGVycyI6W3sibmFtZSI6InBhcmFtMCIsInNjaGVtYSI6eyIkcmVmIjoiIy9jb21wb25lbnRzL3NjaGVtYXMvTWVzc2FnZSJ9fSx7Im5hbWUiOiJwYXJhbTEiLCJzY2hlbWEiOnsidHlwZSI6InN0cmluZyJ9fV0sInRhZyI6WyJzdWJtaXQiLCJTVUJNSVQiXSwibmFtZSI6IlB1dFZhbHVlIiwicmV0dXJucyI6eyJ0eXBlIjoic3RyaW5nIn19LHsidGFnIjpbInN1Ym1pdCIsIlNVQk1JVCJdLCJuYW1lIjoiUmVub3VuY2VPd25lcnNoaXAifSx7InBhcmFtZXRlcnMiOlt7Im5hbWUiOiJwYXJhbTAiLCJzY2hlbWEiOnsiJHJlZiI6IiMvY29tcG9uZW50cy9zY2hlbWFzL01lc3NhZ2UifX0seyJuYW1lIjoicGFyYW0xIiwic2NoZW1hIjp7InR5cGUiOiJhcnJheSIsIml0ZW1zIjp7InR5cGUiOiJpbnRlZ2VyIiwiZm9ybWF0IjoiaW50MzIiLCJtYXhpbXVtIjoyNTUsIm1pbmltdW0iOjB9fX0seyJuYW1lIjoicGFyYW0yIiwic2NoZW1hIjp7InR5cGUiOiJzdHJpbmcifX1dLCJ0YWciOlsic3VibWl0IiwiU1VCTUlUIl0sIm5hbWUiOiJSZW5vdW5jZVJvbGUifSx7InBhcmFtZXRlcnMiOlt7Im5hbWUiOiJwYXJhbTAiLCJzY2hlbWEiOnsidHlwZSI6ImFycmF5IiwiaXRlbXMiOnsidHlwZSI6ImludGVnZXIiLCJmb3JtYXQiOiJpbnQzMiIsIm1heGltdW0iOjI1NSwibWluaW11bSI6MH19fSx7Im5hbWUiOiJwYXJhbTEiLCJzY2hlbWEiOnsidHlwZSI6InN0cmluZyJ9fV0sInRhZyI6WyJzdWJtaXQiLCJTVUJNSVQiXSwibmFtZSI6IlJldm9rZVJvbGUifSx7InBhcmFtZXRlcnMiOlt7Im5hbWUiOiJwYXJhbTAiLCJzY2hlbWEiOnsidHlwZSI6ImFycmF5IiwiaXRlbXMiOnsidHlwZSI6ImludGVnZXIiLCJmb3JtYXQiOiJpbnQzMiIsIm1heGltdW0iOjI1NSwibWluaW11bSI6MH19fSx7Im5hbWUiOiJwYXJhbTEiLCJzY2hlbWEiOnsidHlwZSI6ImFycmF5IiwiaXRlbXMiOnsidHlwZSI6ImludGVnZXIiLCJmb3JtYXQiOiJpbnQzMiIsIm1heGltdW0iOjI1NSwibWluaW11bSI6MH19fV0sInRhZyI6WyJzdWJtaXQiLCJTVUJNSVQiXSwibmFtZSI6IlNldFJvbGVBZG1pbiJ9LHsidGFnIjpbInN1Ym1pdCIsIlNVQk1JVCJdLCJuYW1lIjoiVG90YWwiLCJyZXR1cm5zIjp7InR5cGUiOiJudW1iZXIiLCJmb3JtYXQiOiJkb3VibGUiLCJtYXhpbXVtIjoxODQ0Njc0NDA3MzcwOTU1MjAwMCwibWluaW11bSI6MCwibXVsdGlwbGVPZiI6MX19LHsicGFyYW1ldGVycyI6W3sibmFtZSI6InBhcmFtMCIsInNjaGVtYSI6eyJ0eXBlIjoic3RyaW5nIn19XSwidGFnIjpbInN1Ym1pdCIsIlNVQk1JVCJdLCJuYW1lIjoiVHJhbnNmZXJPd25lcnNoaXAifV0sImRlZmF1bHQiOnRydWV9LCJvcmcuaHlwZXJsZWRnZXIuZmFicmljIjp7ImluZm8iOnsidGl0bGUiOiJvcmcuaHlwZXJsZWRnZXIuZmFicmljIiwidmVyc2lvbiI6ImxhdGVzdCJ9LCJuYW1lIjoib3JnLmh5cGVybGVkZ2VyLmZhYnJpYyIsInRyYW5zYWN0aW9ucyI6W3sidGFnIjpbImV2YWx1YXRlIiwiRVZBTFVBVEUiXSwibmFtZSI6IkdldE1ldGFkYXRhIiwicmV0dXJucyI6eyJ0eXBlIjoic3RyaW5nIn19XSwiZGVmYXVsdCI6ZmFsc2V9fSwiY29tcG9uZW50cyI6eyJzY2hlbWFzIjp7Ik1lc3NhZ2UiOnsiJGlkIjoiTWVzc2FnZSIsInByb3BlcnRpZXMiOnsibm9uY2UiOnsidHlwZSI6Im51bWJlciIsImZvcm1hdCI6ImRvdWJsZSIsIm1heGltdW0iOjE4NDQ2NzQ0MDczNzA5NTUyMDAwLCJtaW5pbXVtIjowLCJtdWx0aXBsZU9mIjoxfSwicHVibGljS2V5Ijp7InR5cGUiOiJhcnJheSIsIml0ZW1zIjp7InR5cGUiOiJpbnRlZ2VyIiwiZm9ybWF0IjoiaW50MzIiLCJtYXhpbXVtIjoyNTUsIm1pbmltdW0iOjB9fSwic2lnbmF0dXJlIjp7InR5cGUiOiJhcnJheSIsIml0ZW1zIjp7InR5cGUiOiJpbnRlZ2VyIiwiZm9ybWF0IjoiaW50MzIiLCJtYXhpbXVtIjoyNTUsIm1pbmltdW0iOjB9fX0sInJlcXVpcmVkIjpbIm5vbmNlIiwicHVibGljS2V5Iiwic2lnbmF0dXJlIl0sImFkZGl0aW9uYWxQcm9wZXJ0aWVzIjpmYWxzZX19fX0="
}
```

> Note: `content` is a base64 encoded json

### GET /basic/total

Used to get total count of depositories

#### Example

```shell
curl -X GET \
  http://localhost:9999/basic/total
```

#### Response

```
{
    "total": 0
}
```

### Get /basic/getValue

Used to get depository value

#### Example

We support to get depository value by its `index` or `kid`

1. By index

```shell
curl -X GET \
  'http://localhost:9999/basic/getValue?index=1'
```

2. By kid(key id)

```shell
curl -X GET \
  'http://localhost:9999/basic/getValue?kid=xxx'
```

#### Response

```json
{
  "index": "1",
  "kid": "xxx",
  "value": "xxx"
}
```

### POST /basic/putValue

Used to create a depository with value

#### Example

```shell
curl -X POST \
  http://localhost:9999/basic/putValue \
  -H 'content-type: application/json' \
  -d '{ 
 "message": "base64_encoded_string_of_message",
 "value": "xxx"
}'
```

#### Response

```json
{
  "kid": "xxxxx"
}
```

### POST /basic/verifyValue

Used to verify a depository with value

#### Example

```shell
curl -X POST \
  http://localhost:9999/basic/verifyValue \
  -H 'content-type: application/json' \
  -d '{ 
 "kid": "xxx",
  "index": "xxx",
 "value": "xxx"
}'
```

#### Response

```json
{
  "status": "xxxxx",
  "reason": "xxx"
}
```

### GET /basic/depositories

List depositories

```shell
curl http://localhost:9999/basic/depositories
```

`query参数`:
| name | description | required | default |
| :--: | :--: | :--: | :--: |
| from | pagination | N | 0 |
| size | pagination | N | 10 |
| startTime | start time (1234), unit is second | N | 0 |
| endTime | end time  | N | 0 |
| name | depository name | N | |
| contentName | file name or some description | N | |
| kid | depository id | N | |

```json
{"count":4,"data":[{"index":"25","kid":"5651d9ae0e5a834afda3fac0e1e743ff3ced5e9d","platform":"bestchains","operator":"","owner":"","blockNumber":42,"name":"abc","contentName":"file name","contentID":"some hash","contentType":"some hash","trustedTimestamp":"1682406287"}]}
```

### GET /basic/depositories/:kid

Get depository by kid

```shell
curl http://localhost:9999/basic/depositories/5651d9ae0e5a834afda3fac0e1e743ff3ced5e9d
```

```json
{"index":"22","kid":"28fd8a24340220857c9857dbeb3f365e505951ca","platform":"bestchains","operator":"","owner":"","blockNumber":37,"name":"dep1","contentName":"","contentID":"lk7234jjsdfsf","contentType":"lk7234jjsdfsf","trustedTimestamp":"1682405989"}
```

### GET /basic/depositories/certificate/:kid

**Request:**

Get depository certificate by kid

```shell
curl -XGET http://localhost:9999/basic/depositories/certificate/02e853e0f68566e62fddd9c4e014db65b7f315d9
```

**Response:**

A pdf will be downloaded directly from browser with filename `02e853e0f68566e62fddd9c4e014db65b7f315d9.pdf`(`{kid}.pdf`):

- Content-Type: application/octet-stream
- Content-Disposition: attachment; filename=02e853e0f68566e62fddd9c4e014db65b7f315d9.pdf
