# Kafka and Zookeeper
FROM alpine:3.19

# Install required packages
RUN apk add --update openjdk11-jre supervisor bash curl

# Update to newer Zookeeper version with better ARM64 support
ENV ZOOKEEPER_VERSION 3.9.3
ENV ZOOKEEPER_HOME /opt/apache-zookeeper-${ZOOKEEPER_VERSION}

# Download and install Zookeeper - fixed URL for 3.9.3
RUN wget -q https://dlcdn.apache.org/zookeeper/zookeeper-${ZOOKEEPER_VERSION}/apache-zookeeper-${ZOOKEEPER_VERSION}.tar.gz -O /tmp/zookeeper-${ZOOKEEPER_VERSION}.tgz \
    && tar xfz /tmp/zookeeper-${ZOOKEEPER_VERSION}.tgz -C /opt \
    && rm /tmp/zookeeper-${ZOOKEEPER_VERSION}.tgz

# Add Zookeeper configuration
ADD assets/conf/zoo.cfg $ZOOKEEPER_HOME/conf/

# Update to newer Kafka version with better ARM64 support
ENV SCALA_VERSION 2.13
ENV KAFKA_VERSION 3.7.0
ENV KAFKA_HOME /opt/kafka_${SCALA_VERSION}-${KAFKA_VERSION}

# Download and install Kafka - fixed URL format
RUN wget -q https://archive.apache.org/dist/kafka/${KAFKA_VERSION}/kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz -O /tmp/kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz \
    && tar xfz /tmp/kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz -C /opt \
    && rm /tmp/kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz

# Add startup scripts
ADD assets/scripts/start-kafka.sh /usr/bin/start-kafka.sh
ADD assets/scripts/start-zookeeper.sh /usr/bin/start-zookeeper.sh

# Make the scripts executable
RUN chmod +x /usr/bin/start-kafka.sh /usr/bin/start-zookeeper.sh

# Supervisor config
ADD assets/supervisor/kafka.ini assets/supervisor/zookeeper.ini /etc/supervisor.d/

# 2181 is zookeeper, 9092 is kafka
EXPOSE 2181 9092

CMD ["supervisord", "-n"]