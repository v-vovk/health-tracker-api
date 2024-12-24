
# Health Tracker API

## Overview
This project is a Health Tracker API that allows users to manage food entities and meals while maintaining detailed logs and error handling. It uses Go with a clean and scalable structure.

## Features
- Food entity CRUD operations.
- Logging using `zap`.
- CI pipeline for linting and testing with `golangci-lint`.
- Robust error handling.
- Follows clean code principles with separated concerns (repository, service, handler).
- Docker setup for Postgres and Redis.

## Project Structure
```
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
5. Run linting:
   ```bash
   make lint
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

## Release Strategy

### Steps for Releasing a New Version

1. **Prepare the `develop` Branch**
    - Ensure `develop` is up to date:
      ```bash
      git checkout develop
      git pull origin develop
      ```
    - Run final checks:
      ```bash
      make lint
      make test
      ```
    - Update `README.md` and `CHANGELOG.md` if needed.

2. **Create a Pull Request**
    - Push the latest changes:
      ```bash
      git push origin develop
      ```
    - Open a PR from `develop` to `main`:
        - **Title**: `Release vX.Y.Z`
        - **Description**: Include features, fixes, and testing status.
    - Request reviews and wait for approval.

3. **Merge the PR**
    - Merge the PR into `main` using GitHub.
    - Pull the latest changes locally:
      ```bash
      git checkout main
      git pull origin main
      ```

4. **Tag the Release**
    - Create a tag for the release:
      ```bash
      git tag -a vX.Y.Z -m "Release version X.Y.Z"
      git push origin vX.Y.Z
      ```
    - Update the GitHub Releases section with details of the release.

5. **Deployment**
    - If ready, deploy the release to production.
    - Automate this step using GitHub Actions or a similar CI/CD tool.

6. **Next Development Cycle**
    - Switch back to `develop`:
      ```bash
      git checkout develop
      ```
    - Create a new branch for the next feature or fix:
      ```bash
      git checkout -b feature/next-feature
      ```

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
