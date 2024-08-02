ARG GO_VERSION=1.22.5
ARG NODE_VERSION=20.16

FROM node:${NODE_VERSION}-bookworm AS base

WORKDIR /app
COPY ./web/ ./
RUN npm install && export NODE_ENV=production npm run build

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
