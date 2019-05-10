FROM golang:latest
EXPOSE 5004

RUN go get github.com/codegangsta/negroni
RUN go get github.com/gorilla/mux
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN go get github.com/unrolled/render

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
ENV GOPATH /app
RUN cd /app ; go install goapi
CMD ["/app/bin/goapi"]
