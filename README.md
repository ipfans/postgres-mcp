# Postgres MCP (Research Project)

A Model Control Protocol (MCP) server implementation for PostgreSQL databases. This project provides a simple HTTP interface to interact with PostgreSQL databases through MCP, allowing you to execute read-only queries and explore database resources.

## Features

- MCP-compliant HTTP server
- Read-only SQL query execution
- Database resource listing
- Environment variable support through `.env` files
- Built with Go and Gin web framework

## Prerequisites

- Go 1.23.6 or later
- PostgreSQL database
- Git

## Installation

```bash
git clone https://github.com/ipfans/postgres-mcp.git
cd postgres-mcp
go mod download
```

## Configuration

You can configure the database connection in two ways:

1. Using command-line flags:

```bash
go run cmd/postgres-mcp/main.go -db "postgres://user:password@localhost:5432/dbname"
```

2. Using environment variables:
   - Create a `.env` file in the project root
   - Add your database URL:
     ```
     DATABASE_URL=postgres://user:password@localhost:5432/dbname
     ```

## Usage

1. Start the server:

```bash
go run cmd/postgres-mcp/main.go
```

The server will start on port 8080 by default.

2. Interact with the MCP endpoints:

- List database resources:

  ```
  POST http://localhost:8080/mcp
  Content-Type: application/json

  {
    "type": "function",
    "name": "resources"
  }
  ```

- Execute a read-only query:

  ```
  POST http://localhost:8080/mcp
  Content-Type: application/json

  {
    "type": "function",
    "name": "query",
    "arguments": {
      "query": "SELECT * FROM your_table LIMIT 10"
    }
  }
  ```

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
