FROM golang:1.18-alpine
WORKDIR /app

EXPOSE 9000

COPY . .

RUN go mod download

RUN go build -o hello-server-app ./server

CMD [ "./hello-server-app" ]