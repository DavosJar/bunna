FROM rust:latest AS builder

# Install system dependencies required by rdkafka-sys (bundled librdkafka C library)
RUN apt-get update && apt-get install -y --no-install-recommends \
    cmake \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the entire project (source, manifests, and lock file)
COPY . .

# Build the release binary
RUN cargo build --release

# Stage 2: Minimal runtime image
FROM debian:sid-slim

# Install runtime shared libraries needed by the binary
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libssl3 \
    && rm -rf /var/lib/apt/lists/*

# Copy the compiled binary from the build stage
COPY --from=builder /app/target/release/hardware-monitor-agent /usr/local/bin/hardware-monitor-agent

ENTRYPOINT ["/usr/local/bin/hardware-monitor-agent"]