FROM golang:latest AS builder

RUN go version
ENV GOPATH=/

WORKDIR /usr/src/webApp


COPY ["go.mod", "go.sum", "./"]
RUN go mod download
RUN go install github.com/cosmtrek/air@latest

EXPOSE 3030

COPY lib ./lib
COPY app ./app

WORKDIR /usr/src/webApp/app/cmd
CMD ["air"]