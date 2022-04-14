# Build application
FROM golang:1.17 AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY ./ ./
RUN go mod download
RUN go build -o /out/test ./cmd/test/main.go

# Run server
FROM alpine:3.15.0
WORKDIR /app
COPY --from=build /out/test ./
CMD [ "./test" ]