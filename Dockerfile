FROM tomsteele/veil-evasion
MAINTAINER Tom Steele <tom@stacktitan.com>

ENV VEIL_LISTENER=localhost:4242
ENV VEIL_OUTPUT_DIR=/usr/share/veil-output
ENV SERVER_LISTENER=0.0.0.0:80
ENV ADMIN_USER=admin
ENV ADMIN_PASS=secret

# install node
WORKDIR /tmp
RUN apt-get install -y g++
RUN curl -o node.tar.gz https://nodejs.org/dist/v0.12.5/node-v0.12.5.tar.gz
RUN tar -zxvf node.tar.gz && cd node-v* && ./configure && make && make install

# install go
ENV GOLANG_VERSION 1.4.2
RUN curl -sSL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz \
    | tar -v -C /usr/src -xz
RUN cd /usr/src/go/src && ./make.bash --no-clean 2>&1
ENV PATH /usr/src/go/bin:$PATH
RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
WORKDIR /go

# copy and restore godep
ADD . /go/src/github.com/tomsteele/veil-evasion-api
RUN go get github.com/tools/godep
RUN cd /go/src/github.com/tomsteele/veil-evasion-api && godep restore
RUN go install github.com/tomsteele/veil-evasion-api

# build client
RUN cd /go/src/github.com/tomsteele/veil-evasion-api/client && npm i && npm run-script build
RUN mkdir -p /root/veil-evasion-api && mv /go/src/github.com/tomsteele/veil-evasion-api/client/dist /root/veil-evasion-api/public

# really bad startup script
RUN echo "cd /root/Veil-Evasion && nohup ./Veil-Evasion.py --rpc &" > /root/start.sh
RUN echo "cd /root/veil-evasion-api && sleep 5 && veil-evasion-api" >> /root/start.sh

ENTRYPOINT ["/bin/bash", "/root/start.sh"]
EXPOSE 80
