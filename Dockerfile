FROM golang:1.22.0-alpine

WORKDIR /app

# Copy and download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main.app .

# Expose the port the app runs on 8080
EXPOSE 8080

# Use CMD to run the binary
CMD ["./main.app"]