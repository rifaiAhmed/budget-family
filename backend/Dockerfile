FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /bin/app ./cmd

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /bin/app /app/app
COPY migrations /app/migrations

EXPOSE 8080

CMD ["/app/app"]
