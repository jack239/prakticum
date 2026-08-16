#!/usr/bin/bash
export IMAGE_NAME=confluentinc/cp-kafka:7.6.1
export CONTAINER_ID=$(docker ps -aqf "ancestor=$IMAGE_NAME")
export TOPIC_NAME=$TOPIC_NAME

# Топик для хранения состояния (key-value)
docker exec -it $CONTAINER_ID kafka-topics \
  --create \
  --topic $TOPIC_NAME \
  --bootstrap-server localhost:9092 \
  --partitions 3 \
  --replication-factor 1 \
  --config cleanup.policy=compact \
  --config min.cleanable.dirty.ratio=0.1 \
  --config segment.ms=3600000