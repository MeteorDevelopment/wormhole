FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg ./pkg
COPY main.go ./
RUN go build -o wormhole

CMD [ "./wormhole" ]
