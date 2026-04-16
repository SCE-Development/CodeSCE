# Frontend
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

ARG VITE_API_URL
ARG VITE_WS_URL
ARG VITE_API_KEY
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_WS_URL=$VITE_WS_URL
ENV VITE_API_KEY=$VITE_API_KEY

COPY frontend/package*.json ./

RUN npm cache clean --force
RUN npm install

COPY frontend/ .

RUN npm run build

# Production stage
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]
# Backend Dockerfile (Go/Gin)

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
