FROM golang:1.14

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...

RUN go get github.com/codegangsta/gin

RUN export PATH=$PATH:$(go env GOPATH)/bin

CMD ["make", "dev"]