# Use an official Golang runtime as a parent image for building
FROM golang:1.23 AS build

# Set the working directory in the container
WORKDIR /app

# Set environment variables
ENV GO111MODULE=on

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-extldflags '-static'" -o main ./cmd/api/main.go

# Ensure the executable has the correct permissions
RUN chmod +x main

# Use a minimal Alpine base image for the final build
FROM alpine

# Set the working directory in the container
WORKDIR /app


# Install tzdata for timezone support
RUN apk add --no-cache tzdata

# Set the timezone environment variable and link the timezone file
ENV TZ=America/Sao_Paulo
RUN ln -sf /usr/share/zoneinfo/$TZ /etc/localtime
RUN ls /etc/localtime

# Copy the executable from the build stage
COPY --from=build /app/main .
COPY --from=build /app/cmd/api/.env .

# Ensure the executable has the correct permissions
RUN chmod +rx main

RUN ls -la

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]