FROM golang:1.21.0-alpine3.18 as builder

RUN apk update && apk upgrade && \
  apk --no-cache --update add git make

WORKDIR /app

COPY . .

RUN go mod tidy && \
  go mod download && \
  go build -v -o engine && \
  chmod +x engine

## Distribution
FROM alpine:latest

# Install dependencies
RUN apk update && apk upgrade && \
  apk --no-cache --update add ca-certificates tzdata && \
  mkdir donut

# Install Doppler CLI
RUN wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
  echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
  apk add doppler

WORKDIR /donut

COPY --from=builder /app/engine /donut
