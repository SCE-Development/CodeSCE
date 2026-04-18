# Frontend Dockerfile
# Build Stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY frontend/package*.json .

RUN npm ci

COPY frontend/ .

RUN npm run build

# Production stage
FROM nginx:alpine AS production

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
