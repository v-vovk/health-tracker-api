
# Health Tracker API

The **Health Tracker API** is a backend service for managing food items, meals, and food groups. It provides RESTful endpoints to create, retrieve, update, and delete records while adhering to a clean and maintainable architecture.

## Features

- **CRUD Operations**: Manage food items, meals, and food groups.
- **Pagination**: Query data with `limit` and `offset` parameters.
- **Error Handling**: Consistent JSON-based error responses.
- **Validation**: Input validation using `go-playground/validator`.
- **Structured Logging**: Logging with `zap` for debugging and monitoring.
- **Middleware**: JSON formatting, request logging, and error recovery.
- **AIR**: live hot reload on changes for dev-mode.

---

## Project Structure

```plaintext
.
├── Makefile
├── cmd
│   └── main.go                # Application entry point
├── docker-compose.yml         # Docker services (PostgreSQL, Redis, etc.)
├── go.mod                     # Go module dependencies
├── go.sum                     # Dependency lock file
├── internal                   # Main application code
│   ├── app
│   │   └── food               # Food module
│   │       ├── factory.go     # Factory for initializing the handler
│   │       ├── handler.go     # Handles HTTP interactions
│   │       ├── model.go       # Defines the Food struct
│   │       ├── repository.go  # Database operations for Food
│   │       ├── routes.go      # Routes for Food endpoints
│   │       └── service.go     # Business logic for Food
│   └── infra
│       ├── config             # Configuration management
│       │   └── config.go
│       ├── db                 # Database connection setup
│       │   └── db.go
│       ├── errors             # Custom error handling
│       │   ├── errors.go
│       │   └── http_errors.go
│       ├── logger             # Logging setup using zap
│       │   └── logger.go
│       └── middleware         # HTTP middleware
│           ├── json.go        # JSON response middleware
│           ├── logging.go     # Request logging middleware
│           └── recovery.go    # Error recovery middleware
├── logs
│   └── app.log                # Log output file
├── migrations                 # Database migrations
│   ├── 000001_init_schema.down.sql
│   └── 000001_init_schema.up.sql
├── pkg                        # Reserved for reusable libraries
└── tmp                        # Temporary files (e.g., Air logs)
    ├── air.log
    └── main
```

---

## Setup Instructions

### Prerequisites

- Go 1.23+ installed.
- Docker and Docker Compose installed.
- Air.
- golang-migrate.

### Clone the Repository

```bash
git clone https://github.com/your-repo/health-tracker-api.git
cd health-tracker-api
```

### Environment Variables

Create a `.env` file in the root directory:

```dotenv
# Environment
ENV=DEV

# Application
APP_PORT=1221

# Database
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_HOST=localhost
DB_PORT=5432
```

---

## Run the Application

### Using Docker Compose

1. Start PostgreSQL and Redis services:
   ```bash
   docker compose up -d
   ```

2. Apply database migrations:
   ```bash
   make migrate-up
   ```

3. Start the application:
   ```bash
   make run
   ```
4. Start the application in dev mode:
   ```bash
   make dev
   ```

### Without Docker

1. Start PostgreSQL and Redis locally.
2. Apply database migrations:
   ```bash
   make migrate-up
   ```
3. Run the application:
   ```bash
   go run cmd/main.go
   ```

---

## Available Endpoints

### Health Check

- **GET** `/health`
  - Returns the health status of the API.

### Food Module

- **GET** `/foods`  
  - Query food items with optional `limit` and `offset`.

- **POST** `/foods`  
  - Add a new food item (requires JSON payload).

- **GET** `/foods/{id}`  
  - Retrieve a food item by its ID.

- **PUT** `/foods/{id}`  
  - Update an existing food item by its ID.

- **DELETE** `/foods/{id}`  
  - Delete a food item by its ID.

---

## Development Workflow

### Using Air for Hot Reload

1. Install Air:
   ```bash
   curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
   ```

2. Run the application with hot reload:
   ```bash
   air
   ```

---

## Testing

To be implemented.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
