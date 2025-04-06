FROM alpine:latest

WORKDIR /app

COPY inventoryApp .

CMD ["./inventoryApp"]
