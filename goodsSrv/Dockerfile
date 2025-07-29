FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/config.json .
CMD ["./app"]