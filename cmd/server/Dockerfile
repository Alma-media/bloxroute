FROM golang:latest

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build /app/cmd/main -o bloxroute

ENTRYPOINT ./bloxroute start
