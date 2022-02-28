#!/bin/bash

host="${HOST:-localhost:8080}"

curl -X POST ${host}/application-inventory/application -d \
'{
    "id": 100,
    "name":"app1",
    "description": "Test application",
    "repository": {
      "url": "https://github.com/rromannissen/tackle-testapp.git"
    },
    "extensions": {
      "binary": "io.konveyor.demo:customers-tomcat:0.0.1-SNAPSHOT:war"
    },
    "businessService": "1"
}' | jq -M .

