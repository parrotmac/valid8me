FROM golang:latest

ENV HTTP_PORT 9000
EXPOSE 9000

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./
RUN go build -v ./

RUN mkdir /go/src/app/workdir

ENTRYPOINT ["/go/src/app/app"]
