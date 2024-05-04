# Use the official Golang image as base
FROM golang:latest AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download and install Go modules
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the Go application based on the build argument passed during build time
ARG APP_NAME

RUN CGO_ENABLED=0 GOOS=linux go build -o ${APP_NAME} ./cmd/${APP_NAME}

# Start a new stage from scratch
FROM alpine:3.18

ARG APP_NAME

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the previous stage
COPY --from=build /app/${APP_NAME} ./bin

# Command to run the executable
CMD ["./bin"]
