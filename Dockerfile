FROM golang:alpine AS buildStage
RUN apk --no-cache add ca-certificates
WORKDIR /telegram-bot
COPY go.mod go.sum ./
RUN  go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build

FROM scratch
WORKDIR /app
COPY --from=buildStage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=buildStage /telegram-bot/telegram-bot .
ENTRYPOINT ["/app/telegram-bot"]