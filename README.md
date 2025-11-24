Go gRPC Auth Service

A robust User Authentication service built with Golang, following Clean Architecture principles and communicating via gRPC.

📂 Project Structure

The project follows a modular structure based on Clean Architecture layers to ensure separation of concerns and testability.

GO-GRPC-AUTH/
├── cmd/                    # Command-line commands (serve, migrate, etc.)
├── config/                 # Configuration loading logic
├── infrastructure/         # Infrastructure setup (Database connections, etc.)
├── proto/                  # Protocol Buffer definitions and generated code
└── user/                   # User Domain Module
    ├── delivery/grpc       # gRPC handlers (Presentation layer)
    ├── domain/             # Domain models and interfaces
    ├── repository/         # Database access layer
    ├── transformer/        # Data transformation logic
    └── usecase/            # Business logic layer
├── config.yaml             # Application configuration
└── main.go                 # Application entry point


🛠️ Prerequisites

Go: 1.20+

Protoc: Protocol Buffer Compiler

Protoc Go Plugins:

protoc-gen-go

protoc-gen-go-grpc

🚀 Getting Started

1. Clone and Install Dependencies

git clone <repository-url>
cd GO-GRPC-AUTH
go mod download


2. Configuration

Ensure your config.yaml file is set up correctly in the root directory. This file likely contains your database credentials and server port configurations.

3. Generate Protobuf Files

If you modify the .proto files, regenerate the Go code using the following command:

protoc \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/user.proto


4. Running the Application

Start the Server:

To start the gRPC server, use the serve command defined in your cmd package:

go run main.go serve


Database Migrations:

Based on the project structure (cmd/migrate.go), you likely have a migration command available:

go run main.go migrate


🏗️ Architecture Overview

This project implements Clean Architecture:

Domain: (user/domain) Contains the core business entities and interface definitions. This layer depends on nothing.

Usecase: (user/usecase) Contains the business rules and logic. It orchestrates the flow of data to and from the domain entities.

Repository: (user/repository) Handles data access (SQL, NoSQL, etc.). It implements interfaces defined in the domain layer.

Delivery: (user/delivery/grpc) The transport layer. In this case, it handles gRPC requests and maps them to the usecase layer.

📜 License

MIT
