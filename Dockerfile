FROM alpine:3.4

RUN apk add --no-cache ca-certificates curl

# Copy binary into image
COPY ./cmd/vatcheck/vatcheck-svc /

# Default values
ENV HTTP_BIND 0.0.0.0
ENV HTTP_PORT 8080
ENV HEALTH_BIND 0.0.0.0
ENV HEALTH_PORT 8090

# Expose ports
EXPOSE 8080 8090

HEALTHCHECK --interval=1m --timeout=2s \
  CMD curl -f http://localhost:8080/_health || exit 1


CMD ["--help"]
ENTRYPOINT ["/vatcheck-svc"]
