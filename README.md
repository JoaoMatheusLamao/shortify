# Shortify

Shortify is a simple URL shortener that uses in-memory persistence.

## Description

This project is a URL shortener built with Go and Gin. It allows users to shorten long URLs and store them in memory for quick access.

## Prerequisites

- Docker

## How to Run

1. Clone the repository:

    ```sh
    git clone https://github.com/JoaoMatheusLamao/shortify
    cd shortify
    ```

2. Build the Docker image:

    ```sh
    docker build -t shortify-image .
    ```

3. Run the Docker container:

    ```sh
    docker run -d -p 8080:8080 --name shortify-container shortify-image ./main
    ```

4. The server will be available at `http://localhost:8080`.

## Endpoints

- `GET /healthcheck`: Checks if the server is running.
- `POST /shorten`: Shortens a URL.
- `GET /:shortURL`: Redirects to the original URL.

## Technologies Used

- Go
- Gin
- Docker

## Contribution

Feel free to contribute to the project. Just open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
