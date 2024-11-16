# Step 1: Build the Go application

FROM golang:1.23 AS go-builder

WORKDIR /app

COPY go/ .

RUN apt-get update && apt-get install -y \
    build-essential \
    gcc \
    g++ \
    libc6-dev \
    mecab \
    libmecab-dev \
    mecab-ipadic-utf8 \
    mecab-utils && \
    export CGO_LDFLAGS="$(mecab-config --libs)" && \
    export CGO_CFLAGS="$(mecab-config --cflags)" && \
    rm -rf /var/lib/apt/lists/* && \
    go mod download && \
    go build -o main .


# - * - * - * - * - * - * - * - *  - * - * - * - * - * - * - * - * - * - * - * - * -#

# Step 2: Build the Node.js application

FROM node:18 AS node-builder

WORKDIR /app

COPY node/package.json node/pnpm-lock.yaml ./
RUN npm install -yg pnpm && pnpm install

COPY node/ .
RUN pnpm build

# - * - * - * - * - * - * - * - *  - * - * - * - * - * - * - * - * - * - * - * - * -#

# Step 3: Combine Go and Node.js applications

FROM debian:bullseye-slim

# Install necessary dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    mecab \
    libmecab-dev \
    mecab-ipadic-utf8 \
    mecab-utils \
    nodejs \
    npm \
    curl && \
    rm -rf /var/lib/apt/lists/*
RUN curl -fsSL https://deb.nodesource.com/setup_23.x | bash - && \
    apt-get install -y nodejs

# Copy the Go and Node.js application files from the build stages
COPY --from=go-builder /app/main /app/go-app/
COPY --from=node-builder /app /app/node-app/

# Set up environment variables for Go
# ENV CGO_LDFLAGS=$(mecab-config --libs)
# ENV CGO_CFLAGS=$(mecab-config --cflags)

RUN cd /app/node-app && npm install -yg pnpm && pnpm install

# Expose the ports for both apps
EXPOSE 8080 3000

# Command to run both Go and Node.js applications
CMD ["sh", "-c", "nohup /app/go-app/main & cd /app/node-app && pnpm start"]
