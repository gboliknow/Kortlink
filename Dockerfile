# Use the official Golang image
FROM golang:1.23-alpine

RUN apk add --no-cache make
# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN make build

# Expose the port the Gin server will run on
EXPOSE 8080

# Command to run the application
CMD ["./bin/api"]
