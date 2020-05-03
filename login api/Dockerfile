FROM golang:1.12.1

RUN go get -d -v github.com/gorilla/mux \
    && go get -d -v github.com/dgrijalva/jwt-go \
    && go get -d -v golang.org/x/crypto/bcrypt \
    && go get -d -v gopkg.in/mgo.v2 \
    && go get -d -v gopkg.in/mgo.v2/bson

RUN mkdir /app

ADD . /app/

WORKDIR /app

RUN go build -o main .

CMD ["/app/main"]