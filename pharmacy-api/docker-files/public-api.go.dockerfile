FROM golang:1.25.3

RUN apt-get update && apt-get install -y tzdata

WORKDIR /public-api

ARG PUBLIC_API_MODULE_PATH
ARG PUBLIC_API_SERVICE_PATH

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./public-api/ ./public-api/

COPY ./shared/ ./shared/

COPY .env ./

RUN go build -o main ./public-api

CMD ["/public-api/main"]