# Build the application from the source
FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Run the tests in the container
FROM build AS run-test-stage
RUN go test -v ./...

FROM alpine:latest

RUN apk --update add ca-certificates curl && rm -rf /var/cache/apk/* && apk add --no-cache curl

WORKDIR /app

EXPOSE 8060

COPY --from=build /app/main /app/.env ./

CMD ["./main"]