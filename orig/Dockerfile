FROM golang

WORKDIR /app

COPY . .

RUN go mod init orig
RUN go build -o out .

CMD ["/app/out"]