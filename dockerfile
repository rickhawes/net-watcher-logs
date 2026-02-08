# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

# Stage 2: Run the application in a minimal image
FROM alpine:latest AS runtime
COPY --from=builder /usr/local/bin/app /usr/local/bin/app
RUN apk add --no-cache tzdata
ENV TZ=America/Los_Angeles
EXPOSE 5555
CMD ["app"]
