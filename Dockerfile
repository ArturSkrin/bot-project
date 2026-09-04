FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bot-project .

FROM alpine:3.20
RUN adduser -D -u 10001 appuser
USER appuser
COPY --from=build /bot-project /usr/local/bin/bot-project
ENTRYPOINT ["/usr/local/bin/bot-project"]