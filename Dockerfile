FROM golang:1.15.7-alpine3.12

# Kafka related settings from github.com/blacktop
LABEL maintainer "https://github.com/zhangwei217245"


ENV KAFKA_VERSION 2.6.0
ENV SCALA_VERSION 2.13

LABEL name="kafka" version=${KAFKA_VERSION}

RUN apk add --no-cache openjdk8-jre bash docker coreutils su-exec git
RUN apk add --no-cache -t .build-deps curl ca-certificates jq \
  && mkdir -p /opt \
  && mirror=$(curl --stderr /dev/null https://www.apache.org/dyn/closer.cgi\?as_json\=1 | jq -r '.preferred') \
  && curl -sSL "${mirror}kafka/${KAFKA_VERSION}/kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz" \
  | tar -xzf - -C /opt \
  && mv /opt/kafka_${SCALA_VERSION}-${KAFKA_VERSION} /opt/kafka \
  && adduser -DH -s /sbin/nologin kafka \
  && chown -R kafka: /opt/kafka \
  && rm -rf /tmp/* \
  && apk del --purge .build-deps


RUN go get github.com/Shopify/sarama/tools/kafka-console-consumer
RUN go get github.com/Shopify/sarama/tools/kafka-console-producer

ADD ./app /app

WORKDIR /app

RUN go build

ENV PATH /sbin:/opt/kafka/bin/:$PATH

WORKDIR /opt/kafka

VOLUME ["/tmp/kafka-logs", "/tmp/go_sorting"]

EXPOSE 9092 2181

# COPY config/log4j.properties /opt/kafka/config/
COPY config/server.properties /opt/kafka/config/
COPY config/zookeeper.properties /opt/kafka/config/
#COPY kafka-entrypoint.sh /kafka-entrypoint.sh
COPY scripts /

# ENV KAFKA_ADVERTISED_HOST_NAME=localhost
# ENV KAFKA_CREATE_TOPICS="source:4:1"

# RUN /kafka-entrypoint.sh kafka-server-start.sh config/server.properties

ENTRYPOINT ["/kafka-entrypoint.sh"]
CMD ["kafka-server-start.sh", "config/server.properties"]

# HEALTHCHECK --interval=5s --timeout=2s --retries=5 CMD bin/health.sh