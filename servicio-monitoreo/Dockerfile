# Stage 1: Build the Rust binary
FROM rust:1-bookworm AS builder

# Install system dependencies required by rdkafka-sys
RUN apt-get update && apt-get install -y --no-install-recommends \
    cmake \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the entire project
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
WORKDIR /app
COPY --from=builder /app/target/release/servicio-monitoreo /app/servicio-monitoreo

EXPOSE 9100
CMD ["/app/servicio-monitoreo"]
