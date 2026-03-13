# Build stage
FROM node:20-alpine AS builder

# Install pnpm
RUN npm install -g pnpm

WORKDIR /app

# Copy package files
COPY frontend/package.json frontend/pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy source code
COPY frontend/ .

# Build the application
RUN pnpm run build

# Production stage
FROM nginx:stable-alpine

# Copy the build output from the builder stage
COPY --from=builder /app/build /usr/share/nginx/html

# Copy a custom nginx configuration to handle the proxy
COPY installation/nginx.conf /etc/nginx/conf.d/default.conf

# Expose the standard HTTP port
EXPOSE 80

# Command to run Nginx
CMD ["nginx", "-g", "daemon off;"]
