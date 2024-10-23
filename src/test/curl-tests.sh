#!/bin/bash
URL=http://0.0.0.0:8081
set -ex
curl --fail -I -X GET $URL
echo "-"
curl --fail -I -X GET $URL/health
echo "-"
curl --fail -I -X GET $URL/events
echo "-"
curl --fail -i -X POST $URL/add -H 'Content-Type: application/json' -d '{"service":"infrasvc","event":"deploy-infrasvc-v0.0.2", "eventType":"deploy-qa"}'
echo "-"
curl --fail -i -X POST $URL/add -H 'Content-Type: application/json' -d '{"service":"infrasvc","event":"deploy-infrasvc-v0.0.2", "eventType":"deploy-qa", "Datetime":"2022-07-02 00:00:00"}'
echo "-"
curl --fail -i -X POST $URL/add -H 'Content-Type: application/json' -d '{"service":"infrasvc","event":"deploy-infrasvc-v0.0.3", "eventType":"deploy-qa"}'
echo "-"
curl --fail -i -X POST $URL/add -H 'Content-Type: application/json' -d '{"service":"infra-service","event":"deploy-infra-service-v0.0.1", "eventType":"deploy-qa"}'
echo -e "\nTests Pass"
