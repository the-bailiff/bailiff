# build stage
FROM golang:1.14.3-alpine3.11 as builder
ENV GO111MODULE=on
WORKDIR /bailiff
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/bailiff/.

# final stage
FROM scratch
COPY --from=builder /bailiff/bailiff /usr/local/bin/
EXPOSE 80
ENTRYPOINT ["bailiff"]
