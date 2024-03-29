FROM golang:1.21.0-alpine

WORKDIR /app

COPY . .

RUN apk add g++ && apk add make && apk add git
RUN make clean && make build

EXPOSE 1323
RUN chmod +x bin/main
CMD ["./bin/main"]
