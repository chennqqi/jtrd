FROM kalilinux/kali-linux-docker
MAINTAINER <q@shellpub.com>
RUN apt-get install wordlists john -y && gunzip /usr/share/wordlists/rockyou.txt.gz 

LABEL maintainer "https://github.com/chennqqi"
LABEL malice.plugin.repository = "https://github.com/chennqqi/jtrd.git"

COPY . /go/src/github.com/chennqqi/jtrd
RUN apt-get install git \
                    go \
  && echo "Building hm week password scanner deamon Go binary..." \
  && export GOPATH=/go \
  && mkdir -p /go/src/golang.org/x \
  && cd /go/src/golang.org/x \
  && git clone https://github.com/golang/net \
  && cd /go/src/github.com/chennqqi/hmbd \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/hmbd \
  && rm -rf /go /usr/local/go /usr/lib/go /tmp/* \
  && apk del --purge .build-deps


RUN chown malice -R /malware
WORKDIR /malware

# Add hmb soft 
ADD http://dl.shellpub.com/hmb/latest/hmb-linux-amd64.tgz /malware/hmb.tgz
RUN tar xvf /malware/hmb.tgz -C /malware
RUN ln -s /malware/hmb /bin/hmb

ENTRYPOINT ["jtrd"]
CMD ["--help"]
