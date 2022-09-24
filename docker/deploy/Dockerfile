FROM golang:1.19

WORKDIR /app

COPY . /app

RUN go mod download

RUN go build -o bin/app src/app/main.go

EXPOSE 3000

CMD ["./bin/app"]