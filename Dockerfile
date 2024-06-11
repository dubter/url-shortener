# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /url-shortner
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /url-shortner/app /url-shortner/cmd/main.go

# Run stage
FROM alpine
WORKDIR /url-shortner
COPY templates templates
COPY --from=builder /url-shortner/app .

EXPOSE 8080
ENTRYPOINT [ "./app", "-env", "/etc/url-shortener/.env"]