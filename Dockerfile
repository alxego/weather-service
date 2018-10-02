FROM golang

RUN mkdir $GOPATH/src/app
ADD . $GOPATH/src/app
WORKDIR $GOPATH/src/app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o weather-service .

ENV LISTEN_PORT 80
CMD ["./weather-service"]
