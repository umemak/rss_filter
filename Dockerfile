# Use the official Golang image to create a build artifact.
FROM golang:1.25.0 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod
COPY go.mod go.sum ./

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Tidy and vendor dependencies
RUN go mod tidy
RUN go mod vendor

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/config.html .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]