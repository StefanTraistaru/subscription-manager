FROM golang:latest

# RUN go mod init main
# RUN export GO111MODULE=off
RUN mkdir /app

WORKDIR /go/src/github.com/go-redis/redis/
RUN git clone https://github.com/go-redis/redis.git .
RUN git checkout v6

RUN go get -d -v github.com/dgrijalva/jwt-go \
    && go get -d -v github.com/rs/cors \
    && go get -d -v github.com/gorilla/mux \
    && go get -d -v github.com/twinj/uuid
    # && go get -d -v github.com/go-redis/redis/v7 \


ADD . /app/

WORKDIR /app

RUN go build -o main .

CMD ["/app/main"]