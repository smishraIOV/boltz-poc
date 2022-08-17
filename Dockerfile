# To run issue this commands:
# In your working directory: mv sample-config-ENV.json config.json
# docker build --tag boltz-poc-server .
# docker run docker run boltz-poc-server

FROM golang:1.18-alpine
RUN apk add git
RUN apk add gcc
RUN apk add musl-dev

WORKDIR /server

COPY server ./

RUN go mod download

RUN go mod tidy

RUN go build -o ./boltz-poc-server

EXPOSE 8080

CMD [ "/server/boltz-poc-server" ]
