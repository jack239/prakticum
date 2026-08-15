#!/usr/bin/sh

docker stop $(docker ps -q)

set -e

docker container prune -f
docker run -d --name kafka -p 9092:9092 apache/kafka:4.0.0

docker exec -it kafka /opt/kafka/bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --topic orders
docker exec -it kafka /opt/kafka/bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --topic users.sum
docker exec -it kafka /opt/kafka/bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --topic users.category
docker exec -it kafka /opt/kafka/bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --topic users-sum-group-table