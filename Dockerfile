FROM docker.arvancloud.ir/alpine:latest

WORKDIR /app

ADD ./main /app

CMD ["./main"]
