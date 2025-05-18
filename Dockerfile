# Use an official Golang image as the base
FROM golang:1.24

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code into the container
COPY . ./

# Build the Go application
RUN go build -o content-management-service ./cmd/server

# Expose the service port
EXPOSE ${SERVICE_PORT}

# Run the application
CMD ["./content-management-service"]