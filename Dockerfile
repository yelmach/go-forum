FROM golang:1.22.6-alpine

RUN apk add --no-cache gcc musl-dev bash

WORKDIR /APP

COPY . .

RUN go mod download

RUN go mod tidy 

RUN go build -o main .

LABEL project="forum"
LABEL version = "1"
LABEL team = "ikazbat, yelmach, asaaoud, amazighi, oanass"

CMD ["./main"]