FROM golang

WORKDIR /app

COPY ./obse ./src
COPY ./data ./data

RUN go mod init obse
RUN go build -o out ./src

CMD ["/app/out"]