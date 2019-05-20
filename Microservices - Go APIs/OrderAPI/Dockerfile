FROM golang:latest
EXPOSE 5003

RUN go get github.com/codegangsta/negroni
RUN go get github.com/gorilla/mux
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN go get github.com/unrolled/render

WORKDIR /go/src/app
COPY . .
RUN go install -v
CMD ["app"]
