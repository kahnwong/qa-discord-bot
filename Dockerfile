FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o /qa-discord-bot

# hadolint ignore=DL3007
FROM gcr.io/distroless/static-debian13:latest AS deploy
COPY --from=build /qa-discord-bot /

CMD ["/qa-discord-bot"]
