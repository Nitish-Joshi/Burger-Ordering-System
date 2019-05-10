FROM golang:latest
EXPOSE 3000

RUN go get github.com/codegangsta/negroni
RUN go get github.com/gorilla/mux


RUN go get github.com/unrolled/render

RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson

WORKDIR /go/src/app
COPY . .
RUN go install -v
CMD ["app"]