FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN go build -o main main.go

FROM golang:alpine as template

WORKDIR /app

COPY . .

CMD go run tools/generate.go

FROM alpine:latest as runner

COPY --from=builder /app/main /app/reverse-proxy
COPY --from=template /app/template /app/template

CMD ["./reverse-proxy"]

