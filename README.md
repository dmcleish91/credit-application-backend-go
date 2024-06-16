# Credit Application Backend (Go)

![Go Version](https://img.shields.io/github/go-mod/go-version/dmcleish91/credit-application-backend-go)
![Issues](https://img.shields.io/github/issues/dmcleish91/credit-application-backend-go)

## Introduction

This project is a backend service for handling credit applications, written in Go. It provides a set of RESTful APIs for managing the credit application process, including submission, and pdf uploads.

### Prerequisites

- Go 1.16 or later
- A running instance of a database (e.g., PostgreSQL)
- [Docker](https://www.docker.com/) (optional, for containerized deployment)

### Steps

1. Clone the repository:
    ```sh
    git clone https://github.com/dmcleish91/credit-application-backend-go.git
    cd credit-application-backend-go
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Configure the application:
   - Copy `.env.example` to `.env` and fill in your configuration details.

4. Build the application:
    ```sh
    go build
    ```
