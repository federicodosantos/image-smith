# Build the application from the source
FROM golang:1.24-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM alpine:latest AS build-release-stage

WORKDIR /app

EXPOSE 8080

COPY --from=build-stage /app/main /app/main

CMD ["./main"]