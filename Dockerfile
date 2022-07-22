FROM golang:1.17

ENV GO111MODULE=on

WORKDIR /var/www

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o app .

CMD ["./app"]