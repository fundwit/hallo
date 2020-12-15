FROM busybox:1.31.1-glibc

COPY hallo /
COPY buildInfo.json /
ENV GIN_MODE release
ENTRYPOINT ["/hallo"]