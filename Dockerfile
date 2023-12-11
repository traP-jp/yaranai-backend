FROM golang:latest

WORKDIR /src
COPY ./src /src

RUN go mod tidy \
  && go build -o ./main 

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64
EXPOSE 8080