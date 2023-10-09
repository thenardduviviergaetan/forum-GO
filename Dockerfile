#syntax=docker/dockerfile:1

FROM golang:1.21.2


LABEL "description"="Docker image of Forum project with SQL database"
LABEL version="1.0"
LABEL author="qdelooze"
LABEL support-contact="delooze.quentin@gmail.com"


WORKDIr /app

ADD . ./
RUN go mod download


RUN go build -o /forum cmd/main.go

CMD ["/forum"]

