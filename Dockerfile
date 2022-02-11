FROM golang:latest

WORKDIR /home


COPY . /home/

RUN cd /home

RUN go build -o library

CMD ["/home/library"]