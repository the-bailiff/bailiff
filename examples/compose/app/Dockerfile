FROM golang:1.14.3-alpine3.11
EXPOSE 80
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
CMD ["app"]
