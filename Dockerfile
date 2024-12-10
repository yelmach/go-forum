FROM golang:1.22.6

WORKDIR /APP

COPY . .

RUN go mod download

RUN go mod tidy 

RUN go build -o main .

LABEL project="forum"
LABEL version = "1"
LABEL team = "ikazbat, yelmach, asaaoud, amazighi, oanass"

CMD ["./main"]