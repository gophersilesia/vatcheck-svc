FROM tianon/true

# Copy binary into image
COPY ./cmd/vatcheck/nx-vatcheck /

# Default values
ENV HTTP_BIND 0.0.0.0
ENV HTTP_PORT 8080
ENV HEALTH_BIND 0.0.0.0
ENV HEALTH_PORT 8090

# Expose port
EXPOSE 8080 8090

CMD ["--help"]
ENTRYPOINT ["/nx-vatcheck"]
