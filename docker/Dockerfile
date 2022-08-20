# build stage
FROM golang:1.18-alpine3.15 AS builder
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 go build --tags netgo -o main --ldflags 'extflags=-static' cmd/main.go

# run stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .
COPY env.yml .

EXPOSE 8080

CMD [ "/app/main" ]
