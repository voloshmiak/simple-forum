# Forum Project

A web application built with Golang, primarily utilizing standard library packages such as `net/http`, `log/slog`, `html/template`, `database/sql`, and others.

## Technologies Used

-   **Language:** Go
-   **Database:** PostgreSQL
-   **Routing:** `net/http`
-   **Templates:** `html/template`
-   **Database Interaction:** `database/sql` with `github.com/jackc/pgx/v5` driver
-   **Migrations:** `github.com/golang-migrate/migrate/v4`
-   **Authentication:** JWT (`github.com/golang-jwt/jwt/v5`)
-   **Environment Variables:** `github.com/ilyakaznacheev/cleanenv`
-   **Password Hashing:** `golang.org/x/crypto/bcrypt`
-   **CSRF Protection:** `github.com/justinas/nosurf`

## How to Run

1.  Ensure you have Go installed and a PostgreSQL database set up.
2.  Create a `.env` file in the project root (or use the existing one) and configure the environment variables (see the "Environment Variables" section).
3.  Execute the command:
    ```bash
    go run cmd/webapp/main.go
    ```
    The application will start at `http://localhost:8070` (or the port specified in `PORT`).

## Migrations

Database migrations are applied automatically when the application starts. Migration files are located in the `pkg/postgres/migrations/` directory.

## Project Structure

```
forum-project
├── .env                  # Environment variables
├── go.mod                # Go dependencies
├── go.sum
├── README.md
├── cmd/
│   └── webapp/           # Application entry point
│       └── main.go
├── internal/             # Internal application logic
│   ├── application/      # Application configuration and assembly
│   ├── auth/             # Authentication logic (JWT)
│   ├── config/           # Application configuration
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── model/            # Data models
│   ├── repository/       # Database interaction logic
│   ├── responder/        # Utilities for HTTP responses (errors)
│   ├── route/            # Route definitions
│   ├── service/          # Business logic
│   └── template/         # HTML template handling
├── pkg/
│   └── postgres/         # PostgreSQL utilities
│       ├── migrations/   # Database migration files
│       └── postgres.go   # PostgreSQL connection logic
└── web/
    ├── static/           # Static files (CSS, JS, images)
    └── templates/        # HTML templates
```

## Environment Variables

The application uses the following environment variables (from the `.env` file):

-   `SERVER_PORT`: Port on which the application will run (default `8070`)
-   `SERVER_READ_TIMEOUT`: Read timeout for HTTP server (default `5s`)
-  `SERVER_WRITE_TIMEOUT`: Write timeout for HTTP server (default `10s`)
-  `SERVER_IDLE_TIMEOUT`: Idle timeout for HTTP server (default `15s`)
-   `DB_HOST`: Database host (default `localhost`)
-   `DB_PORT`: Database port (default `5432`)
-   `DB_NAME`: Database name (default `forum_database`)
-   `DB_USER`: Database user (default `postgres`)
-   `DB_PASSWORD`: Database user password (default empty)
-   `JWT_SECRET`: Secret key for signing JWT tokens (default `some_secret_key`)
-  `JWT_EXPIRATION_HOURS`: JWT token expiration time in hours (default `24`)
-   `APP_ENV`: Application environment, affects template caching (e.g., `development` or `production`, default `development`)
-   `TEMPLATES_PATH`: Path to the HTML templates directory (default `web/templates`)
-   `STATIC_PATH`: Path to the static files directory (default `web/static`)
-   `MIGRATIONS_PATH`: Path to the migrations directory (default `pkg/postgres/migrations`)

## Login Credentials (Examples)

**Administrator:**
-   Email: `admin@email.com`
-   Password: `12345`

You can register other users via the registration form in the application.