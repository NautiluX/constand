FROM golang:alpine AS builder
COPY . /app
WORKDIR /app
RUN go build order.go

FROM alpine:latest
COPY --from=builder /app/order /order
CMD ["/order", "-l"]
