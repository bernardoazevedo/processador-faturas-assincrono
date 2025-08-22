FROM golang:1.24.5

WORKDIR /app

RUN go install github.com/air-verse/air@latest

# RUN go mod download && go mod verify
# RUN go build -o /app/server

ENTRYPOINT ["air"]