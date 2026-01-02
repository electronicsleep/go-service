#!/bin/bash
URL=http://0.0.0.0:8081
set -ex
curl --fail -I -X GET $URL
curl --fail -I -X GET $URL/health
curl --fail -I -X GET $URL/events
curl --fail -i -X POST $URL/add -H 'Content-Type: application/json' -d '{"api_key": "test", "service": "infra", "event": "deploy-infra-v0.0.2", "event_type": "deploy-qa"}'
curl --fail -i -X POST $URL/add -H 'Content-Type: application/json' -d '{"api_key": "test123", "service": "infra", "event": "deploy-infra-v0.0.3", "event_type": "deploy-qa"}'
curl --fail -i -X POST $URL/add -H 'Content-Type: application/json' -d '{"api_key": "test123", "service": "infra", "event": "deploy-infra-v0.0.2", "event_type": "deploy-qa", "datetime": "2022-07-02 00:00:00"}'
curl --fail -i -X POST $URL/add -H 'Content-Type: application/json' -d '{"api_key": "test123", "service": "test", "event": "deploy-test-v0.0.3", "event_type": "deploy-qa"}'
curl --fail -i -X POST $URL/add -H 'Content-Type: application/json' -d '{"api_key": "test123", "service": "service", "event": "deploy-service-v0.0.1", "event_type": "deploy-qa"}'
echo -e "\nTests Pass"
