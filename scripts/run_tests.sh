#!/bin/bash

docker-compose -f docker-compose.test.yml up -d

echo "Waiting for database to be ready..."
sleep 5

go test ./internal/repository -v

docker-compose -f docker-compose.test.yml down 