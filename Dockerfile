FROM golang:1.19.3-buster

WORKDIR /home

COPY . /home

RUN cd /home && go build -o Go-RestAPI

CMD [ "/home/Go-RestAPI" ]