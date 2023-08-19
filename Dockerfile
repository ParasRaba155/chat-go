FROM golang:1.21.0-alpine3.18

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o /api

RUN ls

EXPOSE  8080

ENTRYPOINT [ "/api" ]
