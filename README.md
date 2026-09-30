# High-Performance URL Shortener

A highly scalable, production-ready RESTful API for URL shortening built with Go. This project is engineered with a strong focus on system design, concurrency, and performance optimization to handle high-throughput traffic efficiently.

## Architecture & Core Features

* **Asynchronous Analytics Pipeline:** Leverages Go's native concurrency model (**Goroutines** and **Channels**) to track URL click events in the background. This decouples database writes from the main execution thread, ensuring ultra-low latency (`<2ms`) for redirection operations.
* **Algorithmic Rate Limiting:** Implements a custom **Token Bucket** algorithm from scratch as a middleware. This protects the service against API abuse, traffic bursts, and DDoS attempts without relying on external dependencies.
* **Optimized Database Concurrency:** Configured SQLite with **Write-Ahead Logging (WAL)** mode. This resolves default database locking bottlenecks (`SQLITE_BUSY`) and enables robust, thread-safe concurrent reads and writes under heavy load.
* **Efficient Encoding Algorithm:** Utilizes a custom **Base62** encoding mechanism to generate compact, URL-safe, and collision-resistant short codes.

## Tech Stack

* **Language:** Go (Golang)
* **Framework:** Gin HTTP Web Framework
* **Database:** SQLite (modernc.org/sqlite)

## Getting Started

### Prerequisites
* Go 1.20 or higher installed on your machine.

### Installation & Execution

1. Clone the repository:
   ```bash
   git clone https://github.com/ahmed-elbehidy/url-shortener.git

1. Navigate to the project folder
   ```bash
   cd url-shortener

1. Install dependencies:
   ```bash
   go mod tidy

1. Run the server:
   ```bash
   go run main.go

## API Documentation

### 1. Shorten a URL
* **Endpoint:** `POST /shorten`
* **Content-Type:** `application/json`
* **Request Body:**
  ```json
  {
    "url": "https://www.linkedin.com/in/ahmedelbehidy25/"
  }

* **Success Response (200 OK):**
  ```json
  {
    "original_url": "https://www.linkedin.com/in/ahmedelbehidy25/",
    "short_url": "http://localhost:8080/b"
  }

### 2. Redirect to Original URL
* **Endpoint:** `GET /:code`
* **Description:** Redirects the client to the original long URL and asynchronously increments the click counter in the database.
* **Response:** `302 Found` (Redirect)

---
**Author:** Ahmed El Behidy  
**Contact:** ahmed.behidy@gmail.com  