# build
FROM golang:1.25-alpine AS builder

RUN apk update && apk upgrade && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o main ./cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy config file 
COPY ./cmd/api/config.yaml ./

EXPOSE 8080

CMD [ "./main" ]