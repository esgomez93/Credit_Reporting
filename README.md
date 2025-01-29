```markdown
# Meme as a Service (MaaS)

## Overview

MaaS is a microservice-based API that provides a platform for fetching and serving memes. At its core, it offers a simple yet powerful API endpoint: `GET /memes/`. This API allows clients to retrieve memes based on specified metadata like latitude, longitude, and a free-text query. The service is designed to be scalable, maintainable, and easily extensible.

This project demonstrates a clean architecture approach in Golang, incorporating best practices such as dependency injection, database interactions, and token-based authorization.

## Features

-   **Meme Retrieval:** Fetch memes based on location (latitude, longitude) and a search query.
-   **Token-based Authorization:** Clients are managed through a token system, where each API call consumes tokens.
-   **Real-time Token Balance:** Clients can query their current token balance. (Note: Currently implemented with basic polling; can be enhanced with WebSockets or SSE).
-   **Scalable Architecture:** Designed to handle a large number of requests per second (currently supports 100 RPS, with a roadmap to 10,000 RPS).
-   **Database Integration:** Uses PostgreSQL to store client information, token balances, and API call logs.
-   **Clean Code Structure:** Follows a clean architecture with separate layers for API handling, business logic (service), data access (repository), and database models.
-   **Configuration Management:** Utilizes a YAML configuration file for easy management of server and database settings.

## Project Structure

```

maas/
├── cmd/
│   └── maas/
│       └── main.go        \# Main application entry point
├── pkg/
│   ├── api/
│   │   ├── handler.go   \# HTTP handlers
│   │   └── middleware.go \# Middleware functions
│   ├── service/
│   │   └── meme\_service.go   \# Business logic for memes
│   └── repository/
│       └── meme\_repository.go \# Database interactions
├── internal/
│   ├── store/
│   │   ├── models.go    \# Database models (Client, APICall)
│   │   └── db.go        \# Database connection setup
│   └── config/
│        └── config.go       \# Configuration management
├── utils/
│   └── utils.go         \# Utility functions (e.g., random meme generation)
└── go.mod
└── go.sum
└── config.yaml          \# Configuration file
└── README.md

````

-   **`cmd/`:** Contains the main application entry point.
-   **`pkg/`:**  The core logic of the application, divided into:
    -   **`api/`:**  Handles HTTP requests, responses, and middleware.
    -   **`service/`:** Implements the business logic related to memes and token management.
    -   **`repository/`:**  Handles data access to the database.
-   **`internal/`:**
    -   **`store/`:** Defines database models and database connection setup.
    -   **`config/`:** Manages application configuration.
-   **`utils/`:** Contains utility functions used across the application.

## Getting Started

### Prerequisites

-   Go (version 1.18 or later)
-   PostgreSQL (version 12 or later recommended)
-   Git

### Installation

1.  **Clone the repository:**

    ```bash
    git clone <repository_url>
    cd maas
    ```

2.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

3.  **Database Setup:**
    -   Create a PostgreSQL database (e.g., named `maasdb`).
    -   Update the `config.yaml` file with your database credentials:

    ```yaml
    database:
      host: localhost
      port: 5432
      user: your_db_user
      password: your_db_password
      dbname: maasdb
    ```
    -   Create the tables by using the folowing commands in your database : 
    ```sql
    CREATE TABLE clients (
        client_id SERIAL PRIMARY KEY,
        auth_token TEXT UNIQUE NOT NULL,
        token_balance INTEGER DEFAULT 0
    );

    CREATE TABLE api_calls (
        call_id SERIAL PRIMARY KEY,
        client_id INTEGER REFERENCES clients(client_id),
        timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

    -- Create an index on client_id in the api_calls table for faster lookups
    CREATE INDEX idx_api_calls_client_id ON api_calls (client_id);

    -- Insert a dummy client for testing purposes
    INSERT INTO clients (auth_token, token_balance) VALUES ('test_token', 100);
    ```
### Running the Application

1.  **Start the application:**

    ```bash
    go run cmd/maas/main.go
    ```

    The server will start on port 8000 (or the port specified in `config.yaml`).

### Testing the API

You can use `curl` or a tool like Postman to test the API endpoints:

-   **Get a Meme:**

    ```bash
    curl -H "Authorization: test_token" "http://localhost:8000/memes?lat=40.730610&lon=-73.935242&query=food"
    ```

-   **Add Tokens (replace `test_token` with a valid token and `100` with the desired amount):**

    ```bash
    curl -X POST -H "Authorization: test_token" -H "Content-Type: application/json" -d '{"amount": 100}' http://localhost:8000/addtokens
    ```

-   **Get Token Balance:**

    ```bash
    curl -H "Authorization: test_token" http://localhost:8000/balance
    ```

## API Documentation

### `GET /memes`

Retrieves a meme based on the provided parameters.

**Parameters:**

-   `lat` (float, optional): Latitude of the location.
-   `lon` (float, optional): Longitude of the location.
-   `query` (string, optional): A free-text search query.

**Headers:**

-   `Authorization` (string, required): The client's authentication token.

**Response:**

```json
{
  "meme": "The generated meme text",
  "latitude": "40.730610", // If provided in the request
  "longitude": "-73.935242", // If provided in the request
  "query": "food" // If provided in the request
}
````

**Error Responses:**

  - `401 Unauthorized`: If the `Authorization` header is missing or the token is invalid.
  - `402 Payment Required`: If the client has an insufficient token balance.
  - `500 Internal Server Error`: For any other internal server errors.

### `POST /addtokens`

Adds tokens to a client's balance (typically handled by another service, but this endpoint simulates it).

**Headers:**

  - `Authorization` (string, required): The client's authentication token.
  - `Content-Type`: `application/json`

**Request Body:**

```json
{
    "amount": 150
}
```

**Response:**

  - `200 OK`: If tokens were added successfully.
  - `401 Unauthorized`: If the `Authorization` header is missing.
  - `400 Bad Request`: If the request body is invalid.
  - `500 Internal Server Error`: For any other internal server errors.

### `GET /balance`

Retrieves a client's current token balance.

**Headers:**

  - `Authorization` (string, required): The client's authentication token.

**Response:**

```json
{
  "token_balance": 50
}
```

**Error Responses:**

  - `401 Unauthorized`: If the `Authorization` header is missing or the token is invalid.
  - `500 Internal Server Error`: For any other internal server errors.

## Roadmap to Scaling (10,000 RPS)

The current implementation supports 100 requests per second. Here's a plan to scale it to 10,000 requests per second:

1.  **Database Optimization:**

      - **Indexing:** Add indexes to frequently queried columns (e.g., `auth_token` in the `clients` table).
      - **Connection Pooling:** Ensure efficient database connection pooling to handle a large number of concurrent connections.
      - **Read Replicas:** Implement read replicas to distribute read load (especially for token balance checks).
      - **Caching:** Use a caching layer (e.g., Redis) to cache frequently accessed data like token balances.

2.  **Horizontal Scaling:**

      - **Load Balancing:** Deploy multiple instances of the MaaS application behind a load balancer (e.g., Nginx, HAProxy).
      - **Statelessness:** Ensure the application is stateless so that any instance can handle any request.

3.  **Asynchronous Processing:**

      - **Message Queue:** For operations like logging API calls, use a message queue (e.g., RabbitMQ, Kafka) to offload processing from the main request thread.

4.  **Token Balance Management at Scale:**

      - **Caching:** Aggressively cache token balances in a distributed cache (e.g., Redis) to reduce database load.
      - **Distributed Counters:** Use Redis's atomic increment/decrement operations for efficient token balance updates.
      - **Sharding:** Consider sharding the `clients` table by `client_id` or `auth_token` to distribute data and load across multiple database instances.

5.  **Geographically Diverse Clients:**

      - **CDN:** Use a Content Delivery Network (CDN) to serve static content (if you have image memes) closer to clients.
      - **Multi-Region Deployment:** Deploy the application in multiple regions and use DNS-based routing (e.g., AWS Route 53) to direct clients to the nearest instance.

## CI/CD and Operational SLAs

### CI/CD

  - **Automated Testing:** Implement a comprehensive suite of unit tests, integration tests, and end-to-end tests.
  - **Continuous Integration (CI):** Use a CI platform (e.g., Jenkins, GitLab CI, CircleCI) to automatically build, test, and package the application on every code commit.
  - **Continuous Deployment (CD):** Automate the deployment process to staging and production environments using a CD pipeline.
  - **Blue/Green Deployments:** Use blue/green deployments or rolling updates to minimize downtime during deployments.
  - **Monitoring and Rollbacks:** Implement monitoring to track application health and performance. Automate rollbacks to the previous version in case of deployment failures.

### SLAs

  - **Availability:**
      - Target: 99.95% (approximately 4.38 hours of downtime per year).
      - Measurement: Uptime monitoring using tools like UptimeRobot, Pingdom.
  - **Latency:**
      - Target: 95th percentile latency of under 200ms.
      - Measurement: API performance monitoring using tools like New Relic, Datadog.
  - **Throughput:**
      - Target: 10,000 requests per second (peak).
      - Measurement: Load testing using tools like k6, Gatling.

## Future Enhancements

  - **Meme AI (Premium Feature):** Integrate with a generative AI model to create more unique and dynamic memes. Keep track of client authorization for this feature in the database and cache it for performance.
  - **WebSockets/SSE for Real-time Token Balance:** Implement real-time token balance updates using WebSockets or Server-Sent Events (SSE).
  - **Rate Limiting:** Implement more granular rate limiting per client to prevent abuse.
  - **Advanced Analytics:** Provide clients with dashboards to track their API usage and token consumption.

## Contributing

Contributions to MaaS are welcome\! Please follow these guidelines:

1.  Fork the repository.
2.  Create a new branch for your feature or bug fix.
3.  Write clear and concise code with comments.
4.  Write unit tests for your changes.
5.  Submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](https://www.google.com/url?sa=E&source=gmail&q=LICENSE) file for details.

```
```