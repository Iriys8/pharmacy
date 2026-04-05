FROM golang:1.25.3

RUN apt-get update && apt-get install -y tzdata

WORKDIR /local-api

ARG LOCAL_API_MODULE_PATH
ARG LOCAL_API_SERVICE_PATH

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./local-api/ ./local-api/

COPY ./shared/ ./shared/

COPY .env ./

RUN go build -o main ./local-api

CMD ["/local-api/main"]