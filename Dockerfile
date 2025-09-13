FROM golang:1.23-alpine AS builder

# Install the specific version of the module
ARG APPVERSION="latest"
RUN go install "github.com/eissar/selfhost-crx@${APPVERSION}"

# Copy only the compiled binary to the runtime image
FROM alpine:3.20
# Install certs, create group, create user
RUN apk add --no-cache ca-certificates && update-ca-certificates \
    && addgroup -S appgroup \
    && adduser -S -h /nonexistent -s /sbin/nologin -G appgroup appuser

COPY --from=builder /go/bin/selfhost-crx /usr/local/bin/server
EXPOSE 8080
USER appuser
CMD ["server"]
