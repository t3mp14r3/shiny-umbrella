FROM golang:1.23

WORKDIR /src

COPY . .

RUN go get ./...

RUN go build -o /src/run ./cmd/main.go

CMD ["/src/run"]
