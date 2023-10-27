#syntax=docker/dockerfile:1

FROM golang:1.21.3


LABEL "description"="Docker image of Forum project with SQL database"
LABEL version="1.0"
LABEL author=""
LABEL support-contact=""


WORKDIR /app

ADD . ./
RUN go mod download


# RUN go build -o /admin cmd/usercreation/main.go
RUN go build -o /forum cmd/forum/main.go
# RUN if [-f certgen.sh]; then chmod +x certgen.sh && ./certgen.sh; fi

CMD ["/forum"]

