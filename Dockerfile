FROM golang:1.20-alpine3.17

WORKDIR /app

COPY ms-parking go.mod go.sum ./
COPY config/profiles/default.env ./profiles/

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build .

EXPOSE 80

ENTRYPOINT [ "./ms-parking" ]