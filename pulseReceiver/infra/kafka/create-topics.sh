#!/bin/bash

echo "Criando tópicos no Kafka..."

/usr/bin/kafka-topics \
  --bootstrap-server kafka:9092 \
  --create \
  --if-not-exists \
  --topic aggregated-pulses \
  --partitions 1 \
  --replication-factor 1

echo "Tópico criado com sucesso!"
