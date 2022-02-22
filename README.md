# ATM
Sample code for ATM transactions

Author: Mehran Shokouhi

<br>

## Building the application
```bash
$ go mod tidy
$ go build main.go
$ go test ./...
```

<br>

## Starting the server
```
$ go ./main -config config.json
```

<br>

## Sample config file
```json
{
    "logLevel": "debug",
    "dbPath": "/tmp/atm.db",
    "adminPin": "4444",
    "listenAddr": "localhost:9999"
}
```

<br>

## Overview
This ATM application uses GORM library to connect to Sqlite for persisting its data.

During the startup, the application reads its config from its config file:
- it uses (or creates) the database specified by "dbPath"
- it creates or updates the admin user with userID="admin" and resets its pin to "adminPin"
- it starts an http-server with the address specified by "listenAddr"
- it sets the logging level to "logLevel"

The predefined "admin" can login using the pin in the config file and create users and accounts. After users are created, they can login using their own user-id and pin and work with their accounts. During login process user's pin is compared with the encrypted pin saved in the database. If the pin matches, a session-id will be generated and passed to the user. This session-id together with the user-id must be added to all requests. The session-id is valid until the user logs out.



<br>

## Users
### Admin User

Admin user is a pre-defined user responsible for creating users and account in the system. Admin can login using "admin" and their pin:

```json
POST /session/v1/login HTTP/1.1
Host: localhost:9999
Content-Type: application/json
Content-Length: 45

{
    "user_id": "admin",
    "pin": "4444"
}
```

As the result of the successful login, admin receives a sessionID that must be added to all the request:
```json
{
    "x-user-id": "admin",
    "x-session-id": "ba7531d66dedd15df1ea5e7d02619a38cf573f83490cdd215b5de1f37a6cfb9e",
    "help": "please add 'x-user-id' and 'x-session-id' headers to your requests"
}
```
Without proper x-user-id and x-session-id headers, the request will fail with 401:Unauthorized response.

Admin has access to full functionality of the application including Admin APIs (/admin/v1/.* end-points) and User-APIs (/api/v1/.* end-points).


<br>

### ATM (non-admin) Users
Any user other than admin can only access the User APIs (/api/v1/.* endpoints). ATM users can list their own accounts, deposit/withdraw to/from accounts, or see the list of transactions performed on each account. 


<br>

## APIs

[Postman collection](./docs/atm.postman_collection.json)

### Login API
This is the only API that doesn't need authentication (doesn't need x-user-id and x-session-id headers)

```json
POST /session/v1/login
{
    "user_id": "admin",
    "pin": "4444"
}
```

Response
```json
{
    "x-user-id": "admin",
    "x-session-id": "be47aef167773cc399ef4d3c18b87f2a2dc4b14bfc415dbbadec07dee0c196b4",
    "help": "please add 'x-user-id' and 'x-session-id' headers to your requests"
}
```

<br>

### Admin APIs
Only accessible to "admin" user

#### Create User
```json
POST /admin/v1/users HTTP/1.1
Host: localhost:9999
x-user-id: admin
x-session-id: ba7531d66dedd15df1ea5e7d02619a38cf573f83490cdd215b5de1f37a6cfb9e
Content-Type: application/json
Content-Length: 88

{
    "id": "mike",
    "name": "Mika Mayer",
    "is_active": true,
    "pin": "1234"
}
```

Response:
```json
{
    "created_at": "2022-02-22T15:13:13.766411-05:00",
    "updated_at": "2022-02-22T15:13:13.766411-05:00",
    "deleted_at": "0001-01-01T00:00:00Z",
    "created_by": "admin",
    "updated_by": "",
    "deleted_by": "",
    "id": "mike",
    "name": "Mika Mayer",
    "pin": "a721200c375289479a02127718f6bc25b64c5a5e87a71ec9735f3e1842d34886",
    "is_active": true
}
```

<br>

#### Create Account
```json
POST /admin/v1/accounts HTTP/1.1
Host: localhost:9999
x-user-id: admin
x-session-id: ba7531d66dedd15df1ea5e7d02619a38cf573f83490cdd215b5de1f37a6cfb9e
Content-Type: application/json
Content-Length: 161

{
    "id": "acc-mike-0001",
    "name": "Mike's checking account",
    "currency": "CAD",
    "balance": 125000,
    "owner_id": "mike",
    "is_active": true
}
```

Response
```json
{
    "created_at": "2022-02-22T15:14:57.837448-05:00",
    "updated_at": "2022-02-22T15:14:57.837448-05:00",
    "deleted_at": "0001-01-01T00:00:00Z",
    "created_by": "admin",
    "updated_by": "",
    "deleted_by": "",
    "id": "acc-mike-0001",
    "name": "Mike's checking account",
    "currency": "CAD",
    "balance": 125000,
    "owner_id": "mike",
    "is_active": true
}
```

<bt>

#### List Users
```json
GET /admin/v1/users HTTP/1.1
Host: localhost:9999
x-user-id: admin
x-session-id: ba7531d66dedd15df1ea5e7d02619a38cf573f83490cdd215b5de1f37a6cfb9e
```

Response
```json
[
    {
        "created_at": "2022-02-21T23:17:39.808792-05:00",
        "updated_at": "2022-02-22T14:56:26.546928-05:00",
        "deleted_at": "0001-01-01T00:00:00Z",
        "created_by": "",
        "updated_by": "",
        "deleted_by": "",
        "id": "admin",
        "name": "Admin",
        "pin": "390710ff6385f4961527a32995f372bc663f8152cee7c019f887691a2faa2de2",
        "is_active": true
    },

    ...

    {
        "created_at": "2022-02-22T15:13:13.766411-05:00",
        "updated_at": "2022-02-22T15:13:13.766411-05:00",
        "deleted_at": "0001-01-01T00:00:00Z",
        "created_by": "admin",
        "updated_by": "",
        "deleted_by": "",
        "id": "mike",
        "name": "Mika Mayer",
        "pin": "a721200c375289479a02127718f6bc25b64c5a5e87a71ec9735f3e1842d34886",
        "is_active": true
    }
]
```

<br>

#### List Accounts
```json
GET /admin/v1/accounts HTTP/1.1
Host: localhost:9999
x-user-id: admin
x-session-id: ba7531d66dedd15df1ea5e7d02619a38cf573f83490cdd215b5de1f37a6cfb9e
```

Response
```json
[
    {
        "created_at": "2022-02-21T23:27:58.328771-05:00",
        "updated_at": "2022-02-22T00:12:36.238262-05:00",
        "deleted_at": "0001-01-01T00:00:00Z",
        "created_by": "admin",
        "updated_by": "mehran",
        "deleted_by": "",
        "id": "acc-mehran-0001",
        "name": "Mehran's checking account",
        "currency": "CAD",
        "balance": 150000,
        "owner_id": "mehran",
        "is_active": true
    },
    {
        "created_at": "2022-02-22T15:14:57.837448-05:00",
        "updated_at": "2022-02-22T15:14:57.837448-05:00",
        "deleted_at": "0001-01-01T00:00:00Z",
        "created_by": "admin",
        "updated_by": "",
        "deleted_by": "",
        "id": "acc-mike-0001",
        "name": "Mike's checking account",
        "currency": "CAD",
        "balance": 125000,
        "owner_id": "mike",
        "is_active": true
    }
]
```

<br>

### User APIs
Accessible to all users

#### List Accounts
```json
GET /api/v1/accounts HTTP/1.1
Host: localhost:9999
x-user-id: mike
x-session-id: b89526d67ccbd9d4c75d8a3800a844a1c9f55f9e6f51511a011a4d596fcfa4b7
```

Response
```json
[
    {
        "created_at": "2022-02-22T15:14:57.837448-05:00",
        "updated_at": "2022-02-22T15:14:57.837448-05:00",
        "deleted_at": "0001-01-01T00:00:00Z",
        "created_by": "admin",
        "updated_by": "",
        "deleted_by": "",
        "id": "acc-mike-0001",
        "name": "Mike's checking account",
        "currency": "CAD",
        "balance": 125000,
        "owner_id": "mike",
        "is_active": true
    }
]
```

<br>

#### Get a single Account
```json
GET /api/v1/accounts/acc-mike-0001 HTTP/1.1
Host: localhost:9999
x-user-id: mike
x-session-id: b89526d67ccbd9d4c75d8a3800a844a1c9f55f9e6f51511a011a4d596fcfa4b7
```

Response
```json
{
    "created_at": "2022-02-22T15:14:57.837448-05:00",
    "updated_at": "2022-02-22T15:14:57.837448-05:00",
    "deleted_at": "0001-01-01T00:00:00Z",
    "created_by": "admin",
    "updated_by": "",
    "deleted_by": "",
    "id": "acc-mike-0001",
    "name": "Mike's checking account",
    "currency": "CAD",
    "balance": 125000,
    "owner_id": "mike",
    "is_active": true
}
```

<br>


#### Deposit into Account
```json
POST /api/v1/accounts/acc-mike-0001/transactions HTTP/1.1
Host: localhost:9999
x-user-id: mike
x-session-id: b89526d67ccbd9d4c75d8a3800a844a1c9f55f9e6f51511a011a4d596fcfa4b7
Content-Type: application/json
Content-Length: 45

{
    "type": "deposit",
    "value": 25000
}
```

Response
```json
{
    "created_at": "2022-02-22T15:14:57.837448-05:00",
    "updated_at": "2022-02-22T15:34:30.366089-05:00",
    "deleted_at": "0001-01-01T00:00:00Z",
    "created_by": "admin",
    "updated_by": "mike",
    "deleted_by": "",
    "id": "acc-mike-0001",
    "name": "Mike's checking account",
    "currency": "CAD",
    "balance": 150000,
    "owner_id": "mike",
    "is_active": true
}
```

<br>

#### Withdraw from Account
```json
POST /api/v1/accounts/acc-mike-0001/transactions HTTP/1.1
Host: localhost:9999
x-user-id: mike
x-session-id: b89526d67ccbd9d4c75d8a3800a844a1c9f55f9e6f51511a011a4d596fcfa4b7
Content-Type: application/json
Content-Length: 46

{
    "type": "withdraw",
    "value": 10000
}
```

Response
```json
{
    "created_at": "2022-02-22T15:14:57.837448-05:00",
    "updated_at": "2022-02-22T15:35:55.537693-05:00",
    "deleted_at": "0001-01-01T00:00:00Z",
    "created_by": "admin",
    "updated_by": "mike",
    "deleted_by": "",
    "id": "acc-mike-0001",
    "name": "Mike's checking account",
    "currency": "CAD",
    "balance": 140000,
    "owner_id": "mike",
    "is_active": true
}
```

<br>

#### List Transaction for Account
```
GET /api/v1/accounts/acc-mike-0001/transactions HTTP/1.1
Host: localhost:9999
x-user-id: mike
x-session-id: b89526d67ccbd9d4c75d8a3800a844a1c9f55f9e6f51511a011a4d596fcfa4b7
```

Response
```json
[
    {
        "created_at": "2022-02-22T15:34:30.36715-05:00",
        "updated_at": "2022-02-22T15:34:30.36715-05:00",
        "deleted_at": "0001-01-01T00:00:00Z",
        "created_by": "mike",
        "updated_by": "",
        "deleted_by": "",
        "id": 6,
        "type": "deposit",
        "acc_id": "acc-mike-0001",
        "value": 25000
    },
    {
        "created_at": "2022-02-22T15:35:55.53801-05:00",
        "updated_at": "2022-02-22T15:35:55.53801-05:00",
        "deleted_at": "0001-01-01T00:00:00Z",
        "created_by": "mike",
        "updated_by": "",
        "deleted_by": "",
        "id": 7,
        "type": "withdraw",
        "acc_id": "acc-mike-0001",
        "value": 10000
    }
]
```

<br>


## Some Ideas for Improvements
- Unit tests and Test Coverage
- Extending the functionality (eg. Trasfer between accounts)
- Concurrency 
    
     In the current version, we rely on the DBMS for dealing with access to data. It means if two requests hit the same row in the table, one request can be rejected because of locking. The better way, will be implementing some kind of request-queue, some mechanisms for handling of the requests safely.

- Modularity 
  
    Parts of the code can be moved out as independent libraries to make the code cleaner, easier to read, and easier to test

- Dependency Inversion and Dependency Injection
  
   By replacing the abstract types with interfaces and injecting the dependencies for things like database, session management, encryption..., we can make the code more flexible and improve the testability and maintainability

    ...
