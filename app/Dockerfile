FROM golang:latest AS builder

RUN go version
ENV GOPATH=/

WORKDIR /usr/src/webApp


COPY ["go.mod", "go.sum", "./"]
RUN go mod download
RUN go install github.com/cosmtrek/air@latest

EXPOSE 3030

COPY ./ ./

WORKDIR /usr/src/webApp/cmd
CMD ["air"]