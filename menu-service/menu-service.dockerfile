FROM alpine:latest

WORKDIR /app

COPY menuApp .

CMD ["./menuApp"]
