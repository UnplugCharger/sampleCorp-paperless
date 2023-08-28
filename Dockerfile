# Builder Stage - Build the Go binary
FROM golang:1.20.4-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go


# Final Stage - Copy the binary from the builder stage to the final stage
FROM alpine:3.18.0
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait.sh .
RUN  chmod +x ./wait.sh
COPY db/migrations ./db/migrations
EXPOSE 8090
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]



