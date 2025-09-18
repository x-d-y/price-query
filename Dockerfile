# builder
FROM golang:1.24 AS builder
ARG APP_NAME
WORKDIR /src
COPY . .

RUN if [ "$APP_NAME" = "script" ]; then \
      echo "build script"; \
      make build-history; \
    else \
      echo "build server"; \
      make build; \
    fi

RUN ls -la /src || true
RUN file /src/price-query || true
# final image：Alpine 极简ERROR: failed to build: failed to solve: circular dependency detected on stage: builder

FROM alpine:3.19 AS runtime
COPY --from=builder /src/price-query /app/price-query
RUN chmod +x /app/price-query
ENTRYPOINT ["/app/price-query"]
#ENTRYPOINT ["/bin/sh", "-c", "while true; do sleep 3600; done"]
