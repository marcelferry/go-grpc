FROM golang:1.18-alpine
WORKDIR /app

EXPOSE 50080

ENV GRPC_PORT 50080

COPY . .

RUN go mod download

RUN go build -o hello-server-app ./server

CMD [ "./hello-server-app" ]