# Use the official Golang image with Alpine as the base image.
FROM golang:1.22.6


#Set the working directory inside the container to /app.
WORKDIR /APP


#copy all files in our programme to workdirectory.
COPY . .

#get links in our go.mod and go.sum .
RUN go mod download

#tidy is for delet unessery packges and install the need ones.
RUN go mod tidy 

RUN go build -o main .

#this is the labels for you .
LABEL project="forum"
LABEL version = "1"
LABEL team = ""


#this commond is wil run when a container started from image.
CMD ["./main"]