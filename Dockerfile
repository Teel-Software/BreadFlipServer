FROM golang:1.20

WORKDIR /hleb

COPY . .

ENTRYPOINT make run
