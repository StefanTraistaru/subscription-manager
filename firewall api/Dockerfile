FROM golang:1.12.1

RUN go get -d -v github.com/gorilla/mux \
    && go get -d -v github.com/rs/cors \
    && go get -d -v gopkg.in/mgo.v2 \
    && go get -d -v gopkg.in/mgo.v2/bson \
    && go get -d -v github.com/dgrijalva/jwt-go \
    && go get -d -v github.com/prometheus/client_golang/prometheus \
    && go get -d -v github.com/prometheus/client_golang/prometheus/promhttp

RUN mkdir /app

ADD . /app/

WORKDIR /app

RUN go build -o main .

CMD ["/app/main"]