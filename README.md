
---

# 🚀 Go gRPC Auth Service

A robust **User Authentication service** built with **Golang**, following **Clean Architecture** principles and communicating through **gRPC**.

This project ensures high maintainability, modularity, testability, and clear separation of concerns.

---

## 📁 Project Structure

The project follows a modular folder structure aligned with **Clean Architecture**:

```
GO-GRPC-AUTH/
├── cmd/                 # CLI commands (serve, migrate, etc.)
├── config/              # Configuration loading logic
├── infrastructure/      # Infrastructure setup (DB connections, logger, etc.)
├── proto/               # Protocol Buffers & generated gRPC code
└── user/                # User Domain Module
    ├── delivery/grpc    # gRPC handlers (presentation layer)
    ├── domain/          # Domain models & interfaces (core business rules)
    ├── repository/      # Database access layer (implements domain interfaces)
    ├── transformer/     # Request/response transformation logic
    └── usecase/         # Business logic layer
├── config.yaml          # Application configuration
└── main.go              # Application entry point
```

---

## 🛠️ Prerequisites

Ensure you have the following installed:

* **Go 1.20+**
* **Protocol Buffer Compiler (protoc)**
* **Protoc Go plugins:**

  * `protoc-gen-go`
  * `protoc-gen-go-grpc`

Install plugins:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Make sure `$GOPATH/bin` is in your PATH.

---

## 🚀 Getting Started

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/yourname/GO-GRPC-AUTH.git
cd GO-GRPC-AUTH
go mod download
```

---

### 2️⃣ Configure the Application

Create / update your `config.yaml` file in the root directory.

Example:

```yaml
server:
  port: 50051

database:
  host: localhost
  user: postgres
  password: secret
  name: auth_db
  port: 5432
  sslmode: disable
```

---

### 3️⃣ Generate Protobuf Files

Whenever you modify `.proto` files:

```bash
protoc \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/user.proto
```

---

## ▶️ Running the Application

### Start gRPC Server

```bash
go run main.go serve
```

### Run Database Migrations

```bash
go run main.go migrate
```

---

## 🏗️ Architecture Overview

This project is implemented using **Clean Architecture**, ensuring minimal coupling and maximum flexibility.

### 🧩 Domain Layer (`user/domain`)

* Core business entities (e.g., `User`)
* Domain interfaces
* Contains no external dependencies

### ⚙️ Usecase Layer (`user/usecase`)

* Application-specific business logic
* Coordinates between domain & repository layers

### 🗂 Repository Layer (`user/repository`)

* Database interaction (PostgreSQL, MySQL, etc.)
* Implements domain interfaces

### 🚪 Delivery Layer (`user/delivery/grpc`)

* gRPC request handlers
* Maps incoming gRPC calls to usecases

---

## 📄 License

Distributed under the **MIT License**.
