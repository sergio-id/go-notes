# Step 1: Build the application
FROM golang:1.21.0-alpine3.18 as builder
COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd/user

# Step 3: Add non-root user
FROM alpine:3.18.0 as alpine
RUN addgroup -g 1001 app
RUN adduser app -u 1001 -D -G app /home/app

# Step 4: Final
FROM scratch as final
COPY --from=alpine /etc/passwd /etc/passwd
COPY --from=alpine /etc/group /etc/group
COPY --chown=1001:1001 --from=builder /app/cmd/user/config.yml /config.yml
COPY --chown=1001:1001 --from=builder /bin/app /app

EXPOSE 5004

USER app

CMD ["/app"]
