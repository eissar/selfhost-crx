FROM golang:1.23-alpine AS builder

# Copy local source code and build
WORKDIR /app
COPY . .
RUN go build -o selfhost-crx .

# Copy only the compiled binary to the runtime image
FROM alpine:3.20
# Install certs, create group, create user
RUN apk add --no-cache ca-certificates && update-ca-certificates \
    && addgroup -S appgroup \
    && adduser -S -h /nonexistent -s /sbin/nologin -G appgroup appuser

COPY --from=builder /app/selfhost-crx /usr/local/bin/server
EXPOSE 8080
USER appuser
CMD ["server"]
