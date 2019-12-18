FROM figassis/ubuntu-golang

ENV LOGLEVEL debug
ENV AWS_ACCESS_KEY SOMEKEYID
ENV AWS_SECRET_KEY SOMEKEYID
ENV AWS_REGION us-east-1
ENV BUCKET static.nellcorp.com
ENV PREFIX offload/wp.nellcorp.com
ENV WATCH /data
ENV SCHEDULE "* * * * *"
ENV REDIS_PORT 6379
ENV ACL "private"

COPY . /go/src/github.com/figassis/wpoffload
RUN cd /go/src/github.com/figassis/wpoffload && go mod tidy && go install

VOLUME [ "/data" ]

CMD ["/go/bin/wpoffload"]