FROM golang:1.14.1 AS builder

LABEL maintainer="Loc Ngo <xuanloc0511@gmail.com>"

RUN apt-get update && apt-get install -y apt-utils gcc-aarch64-linux-gnu
# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/goose-linux ./cmd/goose

# Final stage - pick any old arm64 image you want
FROM centos:centos7

RUN mkdir -p /goose-files
WORKDIR /goose-files

COPY --from=builder /app/bin/goose-linux /usr/local/bin/goose
ENTRYPOINT ["/usr/local/bin/goose"]