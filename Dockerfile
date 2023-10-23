FROM golang:1.20

RUN mkdir asdf

RUN mkdir /workok

WORKDIR /hleb

COPY . .

#ENTRYPOINT /bin/bash -c "while true; do sleep 2; done;"
ENTRYPOINT make run
