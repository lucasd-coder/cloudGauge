#!/bin/bash

echo "Criando tópicos no Kafka..."

TOPICS=("aggregated-pulses" "aggregated-to-process")

for topic in "${TOPICS[@]}"; do
  /usr/bin/kafka-topics \
    --bootstrap-server kafka:9092 \
    --create \
    --if-not-exists \
    --topic "$topic" \
    --partitions 1 \
    --replication-factor 1
done

echo "Tópicos criados com sucesso!"
