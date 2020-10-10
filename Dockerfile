FROM golang:1.14

COPY . /go/src/app
WORKDIR /go/src/app
RUN go get
EXPOSE 8080
CMD ["go", "run", "main.go"]