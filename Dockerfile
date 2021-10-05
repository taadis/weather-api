# Use the official golang image to create a build artifact
FROM golang:1.14.15-alpine3.13 as builder

# Create app directory
RUN mkdir /app

# Add file to /app/
ADD . /app/

#
WORKDIR /app

# Update GOPROXY
RUN go env -w GOPROXY=https://goproxy.io,direct && go env GOPROXY

# Build the binary
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o weather-api .

# Run service on container startup
FROM alpine:3.13

WORKDIR /app

COPY --from=builder /app/weather-api .

CMD ["/app/weather-api"]
