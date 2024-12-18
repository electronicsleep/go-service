#!/bin/bash
URL=http://0.0.0.0:8081
set -ex
curl --fail -I -X GET $URL
curl --fail -I -X GET $URL/health
curl --fail -I -X GET $URL/events
curl --fail -i -X POST $URL/add -H 'Content-Type: application/json' -d '{"service":"infra","event":"deploy-infra-v0.0.2", "eventType":"deploy-qa"}'
curl --fail -i -X POST $URL/add -H 'Content-Type: application/json' -d '{"service":"infra","event":"deploy-infra-v0.0.3", "eventType":"deploy-qa"}'
curl --fail -i -X POST $URL/add -H 'Content-Type: application/json' -d '{"service":"infra","event":"deploy-infra-v0.0.2", "eventType":"deploy-qa", "Datetime":"2022-07-02 00:00:00"}'
curl --fail -i -X POST $URL/add -H 'Content-Type: application/json' -d '{"service":"test","event":"deploy-test-v0.0.3", "eventType":"deploy-qa"}'
curl --fail -i -X POST $URL/add -H 'Content-Type: application/json' -d '{"service":"service","event":"deploy-service-v0.0.1", "eventType":"deploy-qa"}'
echo -e "\nTests Pass"
