# Shortify

Shortify is a simple URL shortener that uses in-memory persistence.

## Description

This project is a URL shortener built with Go and Gin. It allows users to shorten long URLs and store them in memory (Redis) for quick access.

Additionally, the project has in-memory IP access control to ensure against hacker attacks (DDoS)

We also save metrics in MongoDB using goroutines to avoid impacting the application's performance.

## Prerequisites

- Docker

## How to Run

1. Clone the repository:

    ```sh
    git clone https://github.com/JoaoMatheusLamao/shortify
    cd shortify
    ```

2. Run the Docker Compose containers:

    ```sh
    docker compose up -d --build
    ```

3. The server will be available at `http://localhost:8080`.

## Endpoints

- `GET /healthcheck`: Checks if the server is running.
- `POST /shorten`: Shortens a URL.

    **Request:**

    ```curl
    curl --location --request POST 'http://localhost:8080/shorten' \
    --header 'Content-Type: application/json' \
    --data '{
        "original_url": "https://youtube.com"
    }'
    ```

    **Response:**

    ```json
    {
        "short_url": "http://localhost:8080/r/12463977768470384210"
    }
    ```

- `GET /r/:shortURL`: Redirects to the original URL.

    **Request:**

    ```curl
    curl --location --request GET 'http://localhost:8080/r/12463977768470384210'
    ```

## Technologies Used

- Go
- Gin
- Docker
- Redis

## Contribution

Feel free to contribute to the project. Just open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
