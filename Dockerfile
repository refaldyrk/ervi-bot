FROM golang:1.23-alpine as builder
LABEL authors="refaldy"

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main .

FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app .

CMD ["./main"]