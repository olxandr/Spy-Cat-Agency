# Spy Cat Agency

This is a sample server for a spy cat agency.

## Prerequisites

- Docker
- Docker Compose

## Setup and Installation

1.  **Clone the repository:**

    ```bash
    git clone git@github.com:olxandr/Spy-Cat-Agency.git
    cd Spy-Cat-Agency
    ```

2.  **Create an environment file:**

    Copy the example environment file to a new `.env` file:

    ```bash
    cp .env.example .env
    ```

    Update the variables in the `.env` file if you need to change the default values.

## Running the Application

To run the application, use Docker Compose:

```bash
docker-compose up --build
```

This will start the following services:

-   **spy-cat-agency:** The main Go application, available at `http://localhost:7777`
-   **db:** A PostgreSQL database, available on port `5432`

To stop the application, run:

```bash
docker-compose down
```

## API Documentation

The API documentation is generated using Swagger and is available at:

[http://localhost:7777/swagger/index.html](http://localhost:7777/swagger/index.html)

## Health Check

The application has a health check endpoint to verify its status:

[http://localhost:7777/healthcheck](http://localhost:7777/healthcheck)

## Environment Variables

The following environment variables are used by the application:

| Variable              | Description                               | Default Value                            |
| --------------------- | ----------------------------------------- | ---------------------------------------- |
| `SPY_CAT_AGENCY_PORT` | The port for the application to run on.   | `:7777`                                  |
| `DB_PROVIDER`         | The database provider.                    | `postgres`                               |
| `DB_HOST`             | The database host.                        | `db`                                     |
| `DB_PORT`             | The database port.                        | `5432`                                   |
| `DB_NAME`             | The database name.                        | `scy`                                    |
| `DB_USER`             | The database user.                        | `cat`                                    |
| `DB_PASSWORD`         | The database password.                    | `scy`                                    |
| `DB_DSN`              | The database connection string.           | `postgres://cat:scy@db/scy?sslmode=disable` |
| `CATS_BREEDS_API`     | The API for fetching cat breeds.          | `https://api.thecatapi.com/v1/breeds`    |
| `DB_MAX_IDLE_TIME`    | The maximum amount of time a connection may be idle. | `15m` |
| `DB_MAX_OPEN_CONNS`   | The maximum number of open connections to the database. | `30` |
| `DB_MAX_IDLE_CONNS`   | The maximum number of connections in the idle connection pool. | `30` |
