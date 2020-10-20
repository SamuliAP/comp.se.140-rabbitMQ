FROM golang

WORKDIR /app

COPY . .

RUN go mod init imed
RUN go build -o out .

CMD ["/app/out"]