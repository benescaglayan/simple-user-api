FROM golang:1.16-alpine

ADD . /go/src/user-service
WORKDIR /go/src/user-service

RUN go build -o /user-service

EXPOSE 8080

ENTRYPOINT [ "/user-service" ]
