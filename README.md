# Mocking Demo Project

This project demonstrates how to use mocking in a microservices architecture with Server A and Server B, along with TestNG automation tests.

## Project Structure

- `server_a/` - Contains the Go implementation of Server A
- `server_b/` - Contains the Go implementation of Server B
- `server_a_tests/` - Contains TestNG automation tests for Server A

## Running the Servers

### Server B

Server B is a gRPC service that processes requests from Server A.

#### QA Environment

In QA mode, Server A uses an internal mock of Server B:
```bash
cd c:\Users\aksha\Desktop\MockingDemo\server_b
go run main.go -env=qa
```

#### stg Environment

In development mode, Server A also uses the internal mock:
```bash
cd c:\Users\aksha\Desktop\MockingDemo\server_b
go run main.go -env=stg
```

Server B will start on port 50051 by default.

### Server A

Server A is a REST API service that forwards requests to Server B. It can be run in different environments.


#### QA Environment

In QA mode, Server A uses an internal mock of Server B:
```bash
cd c:\Users\aksha\Desktop\MockingDemo\server_a
go run main.go -env=qa
```

#### stg Environment

In development mode, Server A also uses the internal mock:
```bash
cd c:\Users\aksha\Desktop\MockingDemo\server_a
go run main.go -env=stg
```

Server A will start on port 8080 by default.

## Running the TestNG Automation Tests

The TestNG tests are designed to test Server A with the internal mock of Server B.

### Prerequisites

- Java 8 or higher
- Maven

### Steps to Run Tests

1. Start Server A in QA mode:
```bash
cd c:\Users\aksha\Desktop\MockingDemo\server_a
go run main.go -env=qa
```

2. In a separate terminal, run the TestNG tests:
```bash
cd c:\Users\aksha\Desktop\MockingDemo\server_a_tests
mvn clean test
```

### Running Specific Tests

To run a specific test class:
```bash
mvn clean test -Dtest=ServerAApiTest
```

To run a specific test method:
```bash
mvn clean test -Dtest=ServerAApiTest#testServerAProcessEndpoint
```

## API Usage

### Server A API

Send a request to Server A:
```bash
curl -X POST http://localhost:8080/api/process -H "Content-Type: application/json" -d "{\"userId\":\"user123\",\"query\":\"test query\"}"
```

## Troubleshooting

- If you encounter connection issues with Server B, make sure it's running before starting Server A in production mode.
- For test failures, check the response body output in the test logs to understand what's being returned from Server A.
