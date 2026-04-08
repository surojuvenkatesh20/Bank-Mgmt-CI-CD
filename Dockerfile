#Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.1/migrate.linux-amd64.tar.gz | tar xvz

#Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY --from=builder /app/migrate /app/migrate
COPY db/migrations /app/migrations
COPY start.sh .
RUN chmod +x start.sh
COPY wait-for.sh .
RUN chmod +x wait-for.sh

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]