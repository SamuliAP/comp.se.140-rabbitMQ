FROM golang

WORKDIR /app

COPY ./httpserv ./src
COPY ./data ./data

RUN go mod init httpserv
RUN go build -o out ./src

CMD ["/app/out"]