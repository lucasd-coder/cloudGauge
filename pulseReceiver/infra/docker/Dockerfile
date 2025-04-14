FROM golang:1.24-alpine AS build

# Instalar gcc e musl-dev
RUN apk add --no-cache gcc musl-dev git wget ca-certificates

RUN addgroup -S appuser \
    && adduser -S -u 10000 -g appuser appuser

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -tags musl -ldflags="-s -w -a" -o /app ./cmd/app

FROM alpine:3.21.3

COPY --from=build /app /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc/passwd /etc/passwd
COPY --from=build /src/config /config

USER appuser

ENTRYPOINT ["/app"]
