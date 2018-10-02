FROM alpine

LABEL author="defectus"
LABEL email="<me@ondralinek.cz>"

RUN apk update && apk add postgresql-client bash busybox ca-certificates libc-utils && addgroup -S glutton && adduser -S -g glutton glutton && mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY glutton-linux-amd64 /usr/local/bin/

USER glutton

CMD /usr/local/bin/glutton-linux-amd64
