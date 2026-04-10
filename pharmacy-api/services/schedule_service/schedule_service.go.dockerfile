FROM golang:1.25.3

WORKDIR /service

COPY ./pharmacy-api .

RUN go mod download

RUN go build -o main ./services/schedule_service

CMD ["/service/main"]