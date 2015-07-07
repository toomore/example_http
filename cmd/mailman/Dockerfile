FROM debian:jessie
MAINTAINER Toomore Chiang <toomore0929@gmail.com>

RUN apt-get update && apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ADD ./mailman /usr/local/bin/

ENV AWSID ""
ENV AWSKEY ""
ENV S3BUCKET ""
ENV S3REGION ""
ENV SESREGION ""
ENV SQSREGION ""
ENV SQSURL ""

ENTRYPOINT mailman
