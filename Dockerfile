FROM golang:1.19

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal

RUN go build -o main cmd/main/main.go
RUN go build -o process cmd/process/main.go
CMD [ "./main" ]
