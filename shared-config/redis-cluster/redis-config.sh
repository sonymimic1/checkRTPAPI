#!/bin/bash
for port in `seq 6371 6376`; do \
  mkdir -p ./shared-config/redis-cluster/redis/${port}/conf \
  && PORT=${port} IP=$1 envsubst < ./shared-config/redis-cluster/redis-cluster.tmpl > ./shared-config/redis-cluster/redis/${port}/conf/redis.conf \
  && mkdir -p ./shared-config/redis-cluster/redis/${port}/data;\
done
