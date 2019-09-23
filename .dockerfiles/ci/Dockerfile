FROM golang:1.12
RUN go get golang.org/x/lint/golint
RUN cd ..
RUN mkdir bollobas
WORKDIR bollobas

ADD . ./
