FROM golang:latest

WORKDIR /src
COPY ./src /src

RUN go mod tidy \
  && go build

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64
EXPOSE 8080

ENV DB_USERNAME="root"
ENV DB_PASSWORD="password"
ENV DB_HOSTNAME="mysqldb"
ENV DB_PORT="3306"
ENV DB_DATABASE="yaranai"
ENV DEVELOPMENT=true