FROM golang:1.25.3

RUN apt-get update && apt-get install -y tzdata

WORKDIR /local-api

COPY ./pharmacy-api/go.mod ./
COPY ./pharmacy-api/go.sum ./

RUN go mod download

COPY ./pharmacy-api/local-api/ ./local-api/

COPY ./pharmacy-api/shared/ ./shared/

RUN go build -o main ./local-api

CMD ["/local-api/main"]