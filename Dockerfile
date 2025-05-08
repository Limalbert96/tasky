# Building the binary of the App
FROM golang:1.19 AS build

WORKDIR /go/src/tasky
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/src/tasky/tasky


FROM alpine:3.17.0 as release

WORKDIR /app
COPY --from=build  /go/src/tasky/tasky .
COPY --from=build  /go/src/tasky/assets ./assets

# Add wizexercise.txt file
RUN echo "This is a deliberately insecure configuration for training purposes" > /app/wizexercise.txt

# Expose port
EXPOSE 8080

# Run as root (intentionally insecure)
USER root

ENTRYPOINT ["/app/tasky"]


