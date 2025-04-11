FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN go build -o main main.go

FROM alpine:latest

COPY --from=builder /app/main /app/main

CMD ["./main"]