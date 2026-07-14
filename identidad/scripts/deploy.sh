#!/bin/bash
cd /home/ubuntu/app
docker-compose -f docker-compose.dev.yml down || true
docker-compose -f docker-compose.dev.yml up -d --build
