FROM golang:1.17.1-alpine3.14 AS builder

# Move to working directory (/build).
WORKDIR /app

# Copy the code into the container.
COPY . .

RUN go build -o main sender/main.go

FROM alpine:3.14
WORKDIR /app
# Copy binary and config files from /build
# to root folder of scratch container.
COPY --from=builder /app/main .

EXPOSE 8080

# Command to run when starting the container.
CMD ["/app/main"]
