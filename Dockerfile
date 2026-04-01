# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy shared dependencies
COPY framework/ ./framework/
COPY libraries/ ./libraries/

# Copy service dependencies
COPY users/go.mod users/go.sum ./users/

WORKDIR /app/users
RUN go mod download

# Copy service code
COPY users/ .

# Build
ARG VERSION=dev
ARG COMMIT=none
ARG BRANCH=unknown
ARG BUILD_TIME=unknown

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X 'yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/version.Version=${VERSION}' \
              -X 'yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/version.Commit=${COMMIT}' \
              -X 'yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/version.Branch=${BRANCH}' \
              -X 'yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/version.BuildTime=${BUILD_TIME}'" \
    -o /go-users main.go

# Stage 2: Run
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies if needed (e.g. ca-certificates, tzdata)
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /go-users ./server
# Copy config if needed, usually passed via env or mounted volume, but copying defaults is safe
COPY --from=builder /app/users/config/config.example.json ./config/config.json

EXPOSE 8080

CMD ["./server", "serve"]
