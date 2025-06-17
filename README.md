# Simple Forum

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
2.  Create a `.env` file in the project root (or use the existing one) and configure the environment variables. See the "Environment Variables" section for details.
3.  Execute the command:
    ```bash
    go run cmd/webapp/main.go
    ```
    The application will start at `http://localhost:8070` (or the port specified in the `SERVER_PORT` environment variable).

## Migrations

Database migrations are applied automatically when the application starts. Migration files are located in the `internal/db/migrations/` directory.

## Project Structure

```
simple-forum
├── .env                  # Environment variables
├── go.mod                # Go dependencies
├── go.sum
├── README.md
├── cmd/
│   └── webapp/           # Application entry point
│       └── main.go
├── internal/             # Internal application logic
│   ├── app/              # Application running configuration and assembly
│   ├── auth/             # Authentication logic (JWT)
│   ├── config/           # Application configuration
│   ├── db/               # Database connection and migration logic
│   │   └── migrations/   # Database migration files
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── model/            # Data models
│   ├── repository/       # Database interaction logic
│   ├── router/           # Route definitions
│   ├── service/          # Business logic
│   └── template/         # HTML template handling
├── web/
    ├── static/           # Static files (CSS, JS, images)
    └── templates/        # HTML templates
```

## Environment Variables

The application uses the following environment variables, which can be configured in the `.env` file. The values from the provided `.env` file are used as examples.

-   `SERVER_PORT`: Port on which the application will run (example: `8070`)
-   `SERVER_READ_TIMEOUT`: Read timeout for HTTP server in seconds (example: `5`)
-   `SERVER_WRITE_TIMEOUT`: Write timeout for HTTP server in seconds (example: `10`)
-   `SERVER_IDLE_TIMEOUT`: Idle timeout for HTTP server in seconds (example: `15`)
-   `DB_ADDR`: Database connection string (example: `postgres://postgres:your_password_here@localhost:5432/forum-database?sslmode=disable`)
-   `JWT_SECRET`: Secret key for signing JWT tokens (example: `your_secret_key_here`)
-   `JWT_EXPIRATION_HOURS`: JWT token expiration time in hours (example: `24`)
-   `APP_ENV`: Application environment, affects template caching (e.g., `development` or `production`, example: `development`)
-   `TEMPLATES_PATH`: Path to the HTML templates directory (example: `web/templates`)
-   `STATIC_PATH`: Path to the static files directory (example: `web/static`)
-   `MIGRATIONS_PATH`: Path to the migrations directory (example: `internal/db/migrations`)

## Login Credentials (Examples)

**Administrator:**
-   Email: `admin@email.com`
-   Password: `12345`

You can register other users via the registration form in the application.