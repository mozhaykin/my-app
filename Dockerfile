FROM golang:1.23.3-alpine3.20 AS build

WORKDIR /app

# Modules layer
COPY go.mod go.sum ./
RUN go mod download

# Build layer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /myapp ./cmd/app

FROM alpine:3.20 AS run

COPY --from=build /myapp /myapp

EXPOSE 8080

CMD ["/myapp"]