FROM google/golang
MAINTAINER Greg Taylor <gtaylor@gc-taylor.com>

WORKDIR /gopath/src/github.com/gtaylor/emdr-relay-go
ADD . /gopath/src/github.com/gtaylor/emdr-relay-go

RUN apt-get update
#RUN apt-cache show libzmq-dev
RUN apt-get install -y mercurial wget build-essential libtool autoconf pkg-config

RUN wget http://download.zeromq.org/zeromq-4.0.5.tar.gz
RUN tar xzf zeromq-4.0.5.tar.gz
WORKDIR zeromq-4.0.5
RUN ./configure && make && make install
RUN ldconfig

WORKDIR ..
RUN go get github.com/pebbe/zmq4

CMD []
EXPOSE 8050
ENTRYPOINT ["go", "run", "emdr-relay-go.go"]
