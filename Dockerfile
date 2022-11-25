FROM golang:1.19

ENV GOPATH=/go/src
ENV WORKSPACE=${GOPATH}/app
RUN mkdir -p ${WORKSPACE}

WORKDIR ${WORKSPACE}

ADD . ${WORKSPACE}

RUN go mod download
RUN go mod tidy

RUN go get -u github.com/cosmtrek/air && go build -o /go/bin/air github.com/cosmtrek/air

CMD ["air", "-c", ".air.toml"]
