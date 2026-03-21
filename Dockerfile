# Backend Dockerfile (Go/Gin)

# --- Dev stage ---
FROM golang:1.25-alpine AS dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 6767

CMD ["go", "run", "./cmd/server"]
# Note: app should listen on port 6767

# --- Build stage ---
FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /bin/server ./cmd/server

# --- Production stage ---
FROM alpine:3.20 AS production

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /bin/server /bin/server

EXPOSE 6767

CMD ["/bin/server"]
