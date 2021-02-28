# Build Caddy v2 with the Authelia module.
FROM caddy:2-builder AS builder
RUN xcaddy build \
    --with github.com/MicahParks/autheliacaddy@v0.0.11

# The actual image being produced with the Authelia module installed.
FROM caddy:2
COPY --from=builder /usr/bin/caddy /usr/bin/caddy
