#Here we use multi-stage build to minimize size of the final image.
#Download go-auto-yt via git and youtube-dl via curl on ubuntu temp image
FROM ubuntu as DOWNLOAD
WORKDIR /git
RUN apt-get update && apt-get install git curl -y && git clone https://github.com/XiovV/go-auto-yt.git && curl -L https://yt-dl.org/downloads/latest/youtube-dl -o ./youtube-dl && chmod a+rx ./youtube-dl

#Transfer git content from DOWNLOAD stage over GO stage to build application
FROM golang:alpine as GO
WORKDIR /app
COPY --from=DOWNLOAD /git/go-auto-yt .
RUN go build -o main .

#Use ffmpeg as base image and copy executable from other temp images
FROM jrottenberg/ffmpeg:alpine as BASE
WORKDIR /app
COPY --from=GO /app/main .
COPY --from=GO /app/static ./static
COPY --from=DOWNLOAD /git/youtube-dl /usr/local/bin/
RUN apk --update add python

#Set starting command
ENTRYPOINT [ "./main" ]

#Expose port and volume
EXPOSE 8080
VOLUME /app/downloads
