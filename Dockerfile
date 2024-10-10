# Stage 1: Build stage
FROM golang:1.23.1-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy all files into the container (Go source code, go.mod, go.sum, etc.)
COPY . .

# Initialize Go modules, tidy dependencies, and build the Go application
RUN go mod tidy && \
    go build -o /go-employee-crud-app

# Stage 2: Run stage (Final image)
FROM alpine:latest

# Set the working directory for the runtime stage
WORKDIR /root/

# Copy the built Go binary from the build stage
COPY --from=builder /go-employee-crud-app ./

# Expose the port where the app will run (if needed, adjust according to your app)
EXPOSE 8080

# Command to run the Go app
CMD ["./go-employee-crud-app"]
