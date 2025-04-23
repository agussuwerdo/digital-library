# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package.json package-lock.json ./

# Install dependencies
RUN npm ci

# Copy the source code
COPY . .

# Create default .env.local if not exists
RUN if [ ! -f .env.local ]; then echo "NEXT_PUBLIC_API_URL=http://localhost:3001" > .env.local; fi

# Build the application
RUN npm run build

# Production stage
FROM node:18-alpine

WORKDIR /app

# Copy built assets from builder
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/next.config.js ./next.config.js
COPY --from=builder /app/.env.local ./.env.local

# Install production dependencies only
RUN npm install --omit=dev

# Expose port 3000
EXPOSE 3000

# Start the application
CMD ["npm", "start"] 