FROM golang:1.25.0-alpine3.22

WORKDIR /spy-cat-agency
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN cp /go/bin/swag /usr/local/bin/swag
RUN apk add --no-cache curl

RUN CGO_ENABLED=0 go build -o /usr/local/bin/api ./cmd/api
RUN apk add --no-cache ca-certificates

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

EXPOSE 7777
ENTRYPOINT ["entrypoint.sh"]
