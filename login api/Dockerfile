FROM golang:1.12.1

# RUN go mod init main
# RUN export GO111MODULE=off

RUN go get -d -v github.com/gorilla/mux \
    # && go get -d -v github.com/mongodb/mongo-go-driver/mongo \
    && go get -d -v go.mongodb.org/mongo-driver/mongo \
    && go get -d -v github.com/dgrijalva/jwt-go \
    && go get -d -v golang.org/x/crypto/bcrypt \
    && go get -d -v gopkg.in/mgo.v2

RUN mkdir /app

ADD . /app/

WORKDIR /app

RUN go build -o main .

CMD ["/app/main"]