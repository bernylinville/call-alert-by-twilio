FROM golang:1.21.0 AS build

WORKDIR /app
COPY . .

# RUN go mod download
# RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/twilio

# FROM debian:12.1-slim

# WORKDIR /app
# COPY --from=build /app/twilio .

# CMD ["./twilio"]

RUN go mod tidy

CMD ["go", "run", "main.go"]
