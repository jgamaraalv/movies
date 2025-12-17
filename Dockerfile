FROM golang:1.24-alpine

WORKDIR /app

# Install git and air
RUN apk add --no-cache git && \
    go install github.com/cosmtrek/air@v1.49.0

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["air", "-c", ".air.toml"]
