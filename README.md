# Kortlink API Documentation

## Overview

Kortlink is a URL shortening service that allows users to create, retrieve, update, and delete short URLs. It also provides statistics on how many times a short URL has been accessed.

The API is built using Go with PostgreSQL and `pgx` for database interactions.

## API Endpoints

### Create Short URL

- **Endpoint:** `POST /shortlink`
- **Description:** Create a new short URL.
- **Request Body:**
  ```json
  {
    "originalURL": "https://example.com"
  }
  ```
- **Response:**
  ```json
  {
    "shortURL": "short.ly/abcd1234"
  }
  ```
- **Errors:**
  - `400 Bad Request`: Invalid request payload or URL format.

### Redirect to Original URL

- **Endpoint:** `GET /:shortURL`
- **Description:** Redirect to the original URL associated with the given short URL.
- **Response:** Redirects to the original URL.

### Update Short URL

- **Endpoint:** `PUT /:shortURL`
- **Description:** Update the original URL for the given short URL.
- **Request Body:**
  ```json
  {
    "originalURL": "https://newexample.com"
  }
  ```
- **Response:**
  ```json
  {
    "message": "Short URL updated successfully"
  }
  ```
- **Errors:**
  - `400 Bad Request`: Invalid request payload or URL format.
  - `404 Not Found`: Short URL does not exist.

### Delete Short URL

- **Endpoint:** `DELETE /:shortURL`
- **Description:** Delete the short URL.
- **Response:**
  ```json
  {
    "message": "Short URL deleted successfully"
  }
  ```
- **Errors:**
  - `404 Not Found`: Short URL does not exist.

### Get Short URL Statistics

- **Endpoint:** `GET /:shortURL/stats`
- **Description:** Retrieve statistics for the given short URL, including the access count.
- **Response:**
  ```json
  {
    "shortURL": "short.ly/abcd1234",
    "originalURL": "https://example.com",
    "accessCount": 42
  }
  ```
- **Errors:**
  - `404 Not Found`: Short URL does not exist.

### Get All Short URLs

- **Endpoint:** `GET /shortlinks`
- **Description:** Retrieve a list of all short URLs and their statistics.
- **Response:**
  ```json
  [
    {
      "shortURL": "short.ly/abcd1234",
      "originalURL": "https://example.com",
      "accessCount": 42
    },
    {
      "shortURL": "short.ly/efgh5678",
      "originalURL": "https://anotherexample.com",
      "accessCount": 15
    }
  ]
  ```

## Store Functions

### CreateShortURL

- **Function:** `CreateShortURL(shortURL *models.ShortURL) error`
- **Description:** Create a new short URL in the database.

### GetOriginalURL

- **Function:** `GetOriginalURL(shortURL string) (string, error)`
- **Description:** Retrieve the original URL for a given short URL.

### IncrementAccessCount

- **Function:** `IncrementAccessCount(shortURL string) error`
- **Description:** Increment the access count for a given short URL.

### UpdateShortURL

- **Function:** `UpdateShortURL(shortURL string, newOriginalURL string) error`
- **Description:** Update the original URL for an existing short URL.

### DeleteShortURL

- **Function:** `DeleteShortURL(shortURL string) error`
- **Description:** Delete a short URL from the database.

### GetShortURLStats

- **Function:** `GetShortURLStats(shortURL string) (*models.ShortURL, error)`
- **Description:** Retrieve statistics for a given short URL.

### GetAllShortURLs

- **Function:** `GetAllShortURLs() ([]models.ShortURL, error)`
- **Description:** Retrieve all short URLs and their statistics.

## Database Setup

The API uses PostgreSQL as the database. The `pgx` library is used for database interactions. Ensure that PostgreSQL is running and the required tables are created.

### Example Docker Command

```bash
docker run -d -p 5432:5432 --name kortlink-db -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=kortlink postgres:latest
```
