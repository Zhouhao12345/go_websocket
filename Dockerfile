FROM golang:1.8

RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/gorilla/websocket
RUN go get github.com/go-redis/redis
RUN go get github.com/satori/go.uuid

WORKDIR /go/src/go_ws
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["go_ws"]
