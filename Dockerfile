# syntax=docker/dockerfile:1.7
FROM golang:1.22 as builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /bin/app /app
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/app"]



