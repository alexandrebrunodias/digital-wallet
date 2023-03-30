FROM golang:1.20

WORKDIR $GOPATH/src/github.com/alexandrebrunodias/wallet-core
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
CMD ["wallet"]