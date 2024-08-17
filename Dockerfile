ARG GO_VERSION=1.22.5
ARG NODE_VERSION=20.16

FROM node:${NODE_VERSION}-bookworm AS base

WORKDIR /app

COPY ./web/package*.json ./
RUN npm ci

COPY ./web/ ./
RUN npm run build:prod

FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
COPY --from=base /app/dist/ ./web/dist/
RUN go build -v -o /run-app .

FROM debian:bookworm
COPY --from=builder /run-app /usr/local/bin/

# アプリケーションを起動
CMD ["run-app"]
EXPOSE 8080
