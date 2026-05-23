# ── build stage ───────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS build
WORKDIR /src

# Cache deps first.
COPY go.mod go.sum ./
RUN go mod download

# Build a static binary (assets/portrait.txt is compiled in via go:embed).
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /sshsite ./...

# ── runtime stage ─────────────────────────────────────────────────────────
FROM alpine:3.20
WORKDIR /app

COPY --from=build /sshsite /app/sshsite

# Fresh production SSH host key, baked in so the server's identity stays
# stable across deploys (no scary "host key changed" warnings for visitors).
COPY .ssh/prod_id_ed25519 /app/.ssh/id_ed25519

EXPOSE 2222
ENTRYPOINT ["/app/sshsite"]
