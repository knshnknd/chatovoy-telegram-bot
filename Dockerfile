FROM golang:1.18.2-alpine3.15

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN mkdir "/app/resources"
COPY resources/ /app/resources

RUN go build -o /chatovoy

CMD ["/chatovoy"]