# We specify the base image we need for our
# go application
FROM golang:1.13.4-alpine3.10
# We update and add git
RUN apk update && apk upgrade
RUN apk add --no-cache git
# We create an /app directory within our
# image that will hold our application source
# files
RUN mkdir /app
# We clone the repository into the app directory
RUN git clone https://github.com/XiovV/go-auto-yt /app
# We specify that we now wish to execute 
# any further commands inside our /app
# directory
WORKDIR /app
# we run go build to compile the binary
# executable of our Go program
RUN go build -o main .
# Our start command which kicks off
# our newly created binary executable
CMD ["/app/main"]