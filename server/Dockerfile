FROM golang:1.24-alpine AS builder


WORKDIR /app
COPY ./src .

RUN go mod tidy && go build -o main .

FROM alpine:latest
WORKDIR /root/

# Установим нужные пакеты
RUN apk add --no-cache tzdata

ENV TZ=Asia/Tashkent

COPY --from=builder /app/main .
COPY ./uploads ./uploads

CMD ["./main"]
