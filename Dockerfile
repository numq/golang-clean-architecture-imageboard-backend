FROM golang:1.17

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["protoc -I . proto/*.proto --go_out=plugins=grpc:.","app"]