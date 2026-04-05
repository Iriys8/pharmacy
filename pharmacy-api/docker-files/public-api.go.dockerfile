FROM golang:1.25.3

RUN apt-get update && apt-get install -y tzdata

WORKDIR /public-api

ARG PUBLIC_API_MODULE_PATH
ARG PUBLIC_API_SERVICE_PATH

COPY ./pharmacy-api/go.mod ./
COPY ./pharmacy-api/go.sum ./

RUN go mod download

COPY ./pharmacy-api/public-api/ ./public-api/

COPY ./pharmacy-api/shared/ ./shared/

RUN go build -o main ./public-api

CMD ["/public-api/main"]