FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
# RUN go mod download
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o main .

FROM alpine:latest
WORKDIR /app
COPY .env* ./
COPY --from=builder /app/main .
CMD ["./main"]