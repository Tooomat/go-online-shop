# 🛒 Go Online Shop API
Proyek ini adalah REST API berbasis Golang menggunakan framework Fiber. API ini mengimplementasikan arsitektur service-repository pattern dengan dukungan:
- Authentication & Authorization menggunakan JWT (user, admin, super_admin).
- Middleware untuk logging (Logrus), tracing request (Trace ID), dan error handling.
- Modular Apps (contoh: product, user, dll).
- Clean Architecture (service tidak bergantung langsung ke Fiber, hanya pakai context.Context).

---

## 🚀 Fitur Utama
- Register & Login User → menghasilkan Access Token (JWT).
- Authorization Middleware → validasi token, inject role ke context.
- CRUD Product
    - Create Product
    - List Product (pagination cursor)
    - Get Product by SKU
- Error Logging & Monitoring → trace dari service → middleware logrus.

---

## DDD arsitekture (Domain Driven Design)
Online Shop sederhana menggunakan bahasa pemrograman Golang dengan pendekatan DDD arsitekture

--- 

## go modul
- [yamlv3](gopkg.in/yaml.v3)
- [sqlx](github.com/jmoiron/sqlx), [driver_mysql](github.com/go-sql-driver/mysql)
- [fiber/v2](github.com/gofiber/fiber/v2)
- [bcrypt](golang.org/x/crypto/bcrypt)
- [jwtv5](github.com/golang-jwt/jwt/v5)
- [testify/require](github.com/stretchr/testify/require)
- [logrus](github.com/sirupsen/logrus)
- [redis](github.com/redis/go-redis/v9)

---

## ▶️ Langkah instalasi
1. Clone repo
    ```bash


    git clone https://github.com/Tooomat/go-online-shop.git
    cd go-online-shop

2. Install dependency
    ```bash

    go mod tidy

3. jalankan
    ```bash

    cd cmd/api
    go run main.go

## ▶️ Langkah instalasi dengan docker
1. jalankan
    ```bash

    docker compose up -d --build

---

## Arsitecture DDD project (*Domain Module Grouping*)
- Tiap domain (misal auth, product, transaction) punya folder sendiri berisi:
entity.go, repository.go, service.go, handler.go, dll.
- Layer seperti handler/service/repo ada di dalam satu domain.

go-online-shop/
├── apps/                       # Domain modules
│   ├── auth/                   # Domain Auth
│   │   ├── entity.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── handler.go
│   ├── product/                # Domain Product
│   └── transaction/            # Domain Transaction
│
├── cmd/                        # Entry point (bootstrap)
│   └── api/
│       └── main.go             # Start HTTP server
│
├── external/                   # External services (DB, cache, third-party)
│   ├── cache/                  # Redis, in-memory cache
│   └── database/               # MySQL/PostgreSQL config & connection
│
├── infrastructure/             # Infrastructure layer
│   ├── http/
│   │   └── fiber/              # HTTP framework adapters (Fiber)
│   ├── middleware/             # Auth middleware, rate limiter, etc.
│   └── response/               # Standardized API responses
│
├── internal/                   # Internal configs (cannot be imported from outside)
│   ├── configs/                # Environment loader
│   └── log/                    # Logger setup
│
├── sql/                        # SQL migrations & queries
│   └── queries/
│
├── test/                       # Unit & integration tests
│
└── utility/                    # Utility global (JWT token, helpers, etc)


---

## Modul authentikasi
    > URL: "/api/v1/auth/regiter" 
        -> method: POST
        -> request body:
            {
                "email": "example@example.com",
                "password": "pass1234"
            }
        -> response: 
               - success
                {
                    "success": true,
                    "message": "register success"
                }
        
               - failed email already used
                {
                    "success": false,
                    "message": "register failed",
                    "error": "email already used",
                    "error_code": "40901"
                }

    > URL: "/api/v1/auth/login"
        -> method: POST
        -> request body:
            {
                "email": "example@example.com",
                "password": "pass1234"
            }
        -> response: 
               - success
                {
                    "success": true,
                    "message": "login success",
                    "payload": {
                        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBBdCI6MTc2MjM1NDk4MSwiaWQiOiIxOTNjYzk5Yi0wYTgzLTQ5YTQtOWNiOC1kNjg3NGE4OTg1NGQiLCJpc3N1ZUF0IjoiMjAyNS0xMS0wNVQxNDo0ODowMS41MzQzMjU0MTFaIiwicm9sZSI6InN1cGVyX2FkbWluIn0.4ZsnfYxCqv7Q-a_FphO4NRkg9rMX5UhDXWXs6lj42d4"
                    }
                }
        
                - failed not found
                {
                    "success": false,
                    "message": "login failed",
                    "error": "not found",
                    "error_code": "40400"
                }

                - failed, password not match
                {
                    "success": false,
                    "message": "login failed",
                    "error": "password not match",
                    "error_code": "40101"
                }

    > URL: "/api/v1/auth/logout"
        -> method: POST
        -> request: Header
                - Authorization: Bearer <token>
        -> response: 
                - success
                {
                    "success": true,
                    "message": "logout success"
                }
        
                - failed login, token missing
                {
                    "success": false,
                    "message": "token missing",
                    "error": "unauthorized",
                    "error_code": "40100"
                }

    > URL: "/api/v1/auth/refresh"
        -> method: POST
        -> response: 
               - success
                {
                    "success": true,
                    "message": "success generate token",
                    "payload": {
                        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBBdCI6MTc2MjM1NDk4MSwiaWQiOiIxOTNjYzk5Yi0wYTgzLTQ5YTQtOWNiOC1kNjg3NGE4OTg1NGQiLCJpc3N1ZUF0IjoiMjAyNS0xMS0wNVQxNDo0ODowMS41MzQzMjU0MTFaIiwicm9sZSI6InN1cGVyX2FkbWluIn0.4ZsnfYxCqv7Q-a_FphO4NRkg9rMX5UhDXWXs6lj42d4"
                    }
                }
        
                - failed
                {
                    "success": false,
                    "message": "failed generate token",
                    "error": "refresh token is missing",
                    "error_code": "40103"
                }
                
---

## Modul product
    > URL:POST "/api/v1/product"
        -> method: POST
        -> request: Header
                - Authorization: Bearer <token>
        -> request body:
            {
                "name": "userabcd",
                "price": 10_000,
                "stock": 20
            }
        -> response:
            - success
            {
                "success": true,
                "message": "create product success"
            }

            - failed create product
            {
                "success": false,
                "message": "create product failed"",
                "error": "product required",
                "error_code": "40005"
            }
    
    > URL: "/api/v1/product"
        -> method: GET
        -> request: Header
                - Authorization: Bearer <token>
        -> request body: cursor and size (default value 0-10)
            {
                "cursor": 0,
                "size": 10
            }
        -> response:
            - success
            {
                "success": true,
                "message": "get list product succcess",
                "payload": [
                    {
                        "id": 1,
                        "sku": "984fad57-5834-489c-bd94-9df4b3bb90d9",
                        "name": "Baju Lebaran",
                        "stock": 20,
                        "price": 30000
                    },
                    {
                        "id": 2,
                        "sku": "7ffedab2-42cb-4d80-96d2-6d9e782fff27",
                        "name": "Baju spiderman",
                        "stock": 12,
                        "price": 16000
                    },
                    {
                        "id": 3,
                        "sku": "cc941c10-db22-431f-a75a-9f380234a79a",
                        "name": "baju poopy",
                        "stock": 9,
                        "price": 17000
                    },
                    ...
                    ...
                ],
                "query": {
                    "cursor": 0,
                    "size": 10
                }
            }

            - success, but empty
            {
                "success": true,
                "message": "get list product succcess",
                "payload": [],
                "query": {
                    "cursor": 9,
                    "size": 8
                }
            }

    > URL: GET "/api/v1/product/sku/:sku"
        -> method: GET
        -> request: Header
            - Authorization: Bearer <token>
        -> response:
            - success
                {
                    "success": true,
                    "message": "get product detail success",
                    "payload": {
                        "id": 3,
                        "sku": "cc941c10-db22-431f-a75a-9f380234a79a",
                        "name": "baju poopy",
                        "stock": 9,
                        "price": 17000,
                        "created_time": "2025-11-10T18:58:01Z",
                        "update_time": "2025-11-10T18:58:01Z"
                    }
                }

            - not found
            {
                "success": false,
                "message": "not found",
                "error": "not found",
                "error_code": "40400"
            }
    
## Modul Transaction
    > URL: POST "api/v1/transactions/checkout"
        -> method: POST
        -> request: 
            - Header
                - Authorization: Bearer <token>
            - body
                {
                    "product_sku": "de209782-bb8f-4e8f-9a71-9d2e2030dbff",
                    "amount": 1
                }
        -> response:
            - success
            {
                "success": true,
                "message": "create transaction success"
            }

            - not found product_sku
            {
                "success": false,
                "message": "transaction failed",
                "error": "not found",
                "error_code": "40400"
            }
            
            - amount 
            {
                "success": false,
                "message": "transaction failed",
                "error": "bad request",
                "error_code": "40000"
            }

    > URL: GET "api/v1/transactions/user/history"
        -> method: GET
        ->	request: Header
                - Authorization: Bearer <token>
        -> response:
            - success
            {
                "success": true,
                "message": "get transaction history success",
                "payload": [
                    {
                        "id": 0,
                        "user_public_id": "917d8170-568a-4076-9cd9-1c78ee2b6933",
                        "product_id": 1,
                        "product_price": 30000,
                        "amount": 2,
                        "sub_total": 60000,
                        "platform_fee": 10000,
                        "grand_total": 70000,
                        "status": "CREATED",
                        "created_time": "2025-11-11T18:40:23Z",
                        "update_time": "2025-11-11T18:40:23Z",
                        "product": {
                            "id": 1,
                            "sku": "de209782-bb8f-4e8f-9a71-9d2e2030dbff",
                            "name": "Baju Lebaran",
                            "price": 30000
                        }
                    }
                ]
            }

---

## Terealisasi
- database pooling, redis
- rate limiter API
- JWT, bcrypt, authentikasi
- authorization, middlerware
- error handling

---

## Soon
- fitur third party/payment gateway {
    use: pembayaran user
    info: untuk payment,bisa dijadiin domain baru. jadi saat dia baru bikin transaction, statusnya diset PENDING. dan saat dia sudah bayar, status kita ubah jadi PAID.
}

## Jika untuk Scalability
- Load Balancer {
    use: untuk mengatur traffic req user ke banyak server database (app 1, app 2, app 3)/node
}
- stateless {
    use: Backend tidak menyimpan session di local memory, melainkan di redis
}
- Database Master-Slave(replication) {
    use: menyediakan server khusus untuk READ (SELECT => MASTER) dan WRITE (INSERT, DELETE, dan UPDATE => SLAVE)   
}
- Sharding (studi kasus) {
    use: dilakukan jika user sudah banyak, atau Membagi data besar ke beberapa database agar scalable
    how: disimpan berdasarkan user_id, (
        Horizontal shard(membagi rows ke beberapa db): tabel user dipecah berdasarkan user_id
        exam -> user_id = 1-10.000 => db-server-1, user_id = 10.000-20.000 => db-server-2
        
        Vertical shard(membagi column ke db berbeda): 
    )
    - jadi 1 shard mempunyai punya 1 MASTER dan beberapa SLAVE
}
- Microservices {

}
- Message Queue/ broker {
    use: gunakan kafka/rabbitMQ untuk mengatur proses berat (email, notif, uploud besar, dll)
}