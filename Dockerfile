FROM golang:1.24.5

WORKDIR /app

COPY . .

RUN go mod download && go mod verify
RUN go build -o /app/server

ENTRYPOINT ["/app/server"]