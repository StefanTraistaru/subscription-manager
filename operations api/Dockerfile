FROM golang:latest

# RUN go get -d github.com/gorilla/mux
RUN go get -d -v github.com/gorilla/mux \
    && go get -d -v github.com/rs/cors \
    && go get -d -v gopkg.in/mgo.v2

RUN mkdir /app

ADD . /app/

WORKDIR /app

RUN go build -o main .

CMD ["/app/main"]