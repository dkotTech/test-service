# Test application with balances and leaderboard

## Running

### From sources

```bash
git clone git@github.com:dkotTech/test-service.git
cd test-service
go run ./cmd/service/main.go 
```

### Docker

```bash
docker run -p 8080:8080 app:latest
```

## Test cases

### WebSocket

```bash
curl 'ws://localhost:8080/ws?event_kind=leaderboard_changes&event_kind=withdraw&event_kind=deposit'
```

### Ok

```bash
curl -X POST 'http://localhost:8080/api/wallet/deposit' \
-d '{"account_id":"00000000-0000-0000-0000-000000000001","amount":10}'

curl -X POST 'http://localhost:8080/api/wallet/withdraw' \
-d '{"account_id":"00000000-0000-0000-0000-000000000001","amount":10}'
```

```bash
curl -X POST 'http://localhost:8080/api/leaderboard/record' \
-d '{"account_id":"00000000-0000-0000-0000-000000000001","score":1}'

curl -X POST 'http://localhost:8080/api/leaderboard/record' \
-d '{"account_id":"00000000-0000-0000-0000-000000000002","score":2}'

curl -X POST 'http://localhost:8080/api/leaderboard/me' \
-d '{"account_id":"00000000-0000-0000-0000-000000000002"}'

curl -X POST 'http://localhost:8080/api/leaderboard/leaders' \
-d '{"limit":100}' 
```

### Bad

1. insufficient funds
```bash
curl -X POST 'http://localhost:8080/api/wallet/deposit' \
-d '{"account_id":"00000000-0000-0000-0000-000000000002","amount":10}'

curl -X POST 'http://localhost:8080/api/wallet/withdraw' \
-d '{"account_id":"00000000-0000-0000-0000-000000000002","amount":100}'
```
```json
{"kind":"user_visible","error":"insufficient funds"}
```

2. withdraw validation
```bash
curl -X POST 'http://localhost:8080/api/wallet/withdraw' \
-d '{"account_id":"7f55f0f8-ebe2-4522-8a96-25a46509885f","amount":-10}'
```
```json
{"kind":"user_visible","error":"Key: 'WithdrawRequest.amount' Error:Field validation for 'amount' failed on the 'gt' tag"}
```

3. deposit validation
```bash
curl -X POST 'http://localhost:8080/api/wallet/deposit' \
-d '{"account_id":"00000000-0000-0000-0000-000000000002","amount":-10}'
```
```json
{"kind":"user_visible","error":"Key: 'DepositRequest.amount' Error:Field validation for 'amount' failed on the 'gt' tag"}
```
