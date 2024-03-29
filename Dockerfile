FROM golang:alpine

RUN apk update && apk add --no-cache git ffmpeg gcc musl-dev

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o binary

ENTRYPOINT ["/app/binary"]
