# Build.
FROM golang:1.14 AS builder

WORKDIR /bot/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o mca_workshop_bot cmd/bot/main.go

# Run.
FROM alpine:3.12.0 AS launcher

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /bot/mca_workshop_bot .

CMD /root/mca_workshop_bot
