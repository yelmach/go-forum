# Use the official Golang image with Alpine as the base image.
FROM golang:1.22.6-alpine

#Set the working directory inside the container to /app.
WORKDIR /APP

#copy all files in our programme to workdirectory.
COPY . .

#get links in our go.mod and go.sum .
RUN go mod download

#tidy is for delet unessery packges and install the need ones.
RUN go mod tidy 

#this is the labels for you .
LABEL project="forum"
LABEL version = "1"/
LABEL team = ""

#here we add bash because alpine doesn't came with it.
RUN apk update && apk add --no-cache bash

#this commond is wil run when a container started from image.
CMD ["go", "run", "."]