FROM golang:1.16
RUN go get github.com/kisielk/errcheck
RUN go get github.com/golang/mock/mockgen
RUN go get github.com/kisielk/godepgraph
RUN apt update
RUN apt install graphviz -y
