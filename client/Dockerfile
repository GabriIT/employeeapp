# 1. Base build stage
FROM node:20 AS builder

WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install

COPY . .
RUN npm run build

# 2. Production runtime image
FROM node:20-alpine

WORKDIR /app
COPY --from=builder /app ./
EXPOSE 3000

ENV NODE_ENV=production
CMD ["npm", "run", "start"]
