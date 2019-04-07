FROM golang

EXPOSE 3000

ADD ./web ${GOPATH}/src/github.com/itouri/fortnite/web

WORKDIR ${GOPATH}/src/github.com/itouri/fortnite/web

RUN go get -u github.com/golang/dep/cmd/dep
RUN make engine
CMD go run server.go
