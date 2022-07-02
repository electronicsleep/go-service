#!/bin/bash
set -ex
curl --fail -I -X GET http://0.0.0.0:8080
echo "-"
curl --fail -I -X GET http://0.0.0.0:8080/health
echo "-"
curl --fail -I -X GET http://0.0.0.0:8080/events
echo "-"
curl --fail -i -X POST http://0.0.0.0:8080/add -H 'Content-Type: application/json' -d '{"service":"infrasvc","event":"deploy-infrasvc-v0.0.2", "eventType":"deploy-qa"}'
echo "-"
curl --fail -i -X POST http://0.0.0.0:8080/add -H 'Content-Type: application/json' -d '{"service":"infrasvc","event":"deploy-infrasvc-v0.0.3", "eventType":"deploy-qa"}'
echo "-"
curl --fail -i -X POST http://0.0.0.0:8080/add -H 'Content-Type: application/json' -d '{"service":"infra-service","event":"deploy-infra-service-v0.0.1", "eventType":"deploy-qa"}'
echo -e "\ntests pass"
