# ATM
sample code for ATM transactions

author: Mehran Shokouhi

## How to build and run
```bash
$ go mod tidy
$ go build main.go
$ go ./main -config config.json
```

## Sample config
```json
{
    "logLevel": "debug",
    "dbPath": "/tmp/atm.db",
    "adminPin": "4444",
    "listenAddr": "localhost:9999"
}
```

## How to test
- Application reads the parameters from the specified config file:
- It uses (or creates) the Sqlite database specified by "dbPath"
- An httpServer will serve at the address specified by "listenAddr"
- The default admin user is "admin" and will be able to login by the pin specified by "adminPin"
- Successful login will provide two headers (x-user-id and x-session-id) that needs to be added to all request to pass the authentication
- "admin" user can create new users and/or accounts
- non-admin users can deposit, withdraw, see their accounts, or transactions of their accounts


## APIs

### Postman Collection
[Postman collection](./docs/atm.postman_collection.json)

### Login API
```
POST /session/v1/login
{
    "user_id": "admin",
    "pin": "4444"
}
```

response
```json
{
    "x-user-id": "admin",
    "x-session-id": "be47aef167773cc399ef4d3c18b87f2a2dc4b14bfc415dbbadec07dee0c196b4",
    "help": "please add 'x-user-id' and 'x-session-id' headers to your requests"
}
```

### User APIs
```
POST|GET /api/v1/.*
x-user-id: mehran
x-session-id: 61f6cf83e1c618e942cd889d52709b45e57ad6155ce88922ffa878b243eae6c2
```

### Admin APIs
```
POST|GET /admin/v1/.*
x-user-id: admin
x-session-id: 709b45e57ad6155ce88922ffa878b243eae6c292f6cf83e1c618e942cd889d52
```





