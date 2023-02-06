#!/bin/bash

set -o errexit
set -o xtrace

host="${HOST:-localhost:8080}"

# Exit early if kubectl or jq not installed
if ! command -v kubectl >/dev/null 2>&1; then
  echo "Please install kubectl"
  exit 1
fi
if ! command -v jq >/dev/null 2>&1; then
  echo "Please install jq"
  exit 1
fi

# Verify we can talk with hub first
if ! curl -S -s -o /dev/null -X GET ${host}/addons/windup; then
  echo "Windup addon not found. Is the hub running?"
  exit 1
fi
echo "Verified windup addon installed."

# Create pathfinder app if it hasn't been added already
# There is a constraint that only allows one application to have a particular name.
if ! curl -S -s -X GET ${host}/applications | jq -e 'any(.[]; .name == "Pathfinder")'; then
  echo "Creating pathfinder application"
  curl -X POST ${host}/applications -d \
    '{
        "name":"Pathfinder",
        "description": "Tackle Pathfinder application.",
        "repository": {
          "name": "tackle-pathfinder",
          "url": "https://github.com/konveyor/tackle-pathfinder.git",
          "branch": "1.2.0"
        }
    }' | jq -M .
fi
APP_ID=$(curl -S -s -X GET ${host}/applications | jq --raw-output '.[] | select(.name=="Pathfinder") | .id')
echo "Pathfinder exists with id ${APP_ID}"
# Show the applications in the inventory
curl -S -s -X GET ${host}/applications | jq

# Make a request to hub
TASK_ID=$(curl -S -s -X POST ${host}/tasks -d \
'{
    "name":"Windup",
    "state": "Ready",
    "locator": "windup",
    "addon": "windup",
    "application": {"id": '$APP_ID'},
    "data": {
        "mode": {
            "artifact": "",
            "binary": false,
            "withDeps": false,
	    "diva": true
        },
        "output": "/windup/report",
        "rules": {
            "path": "",
            "tags": {
                "excluded": [ ]
            }
        },
        "scope": {
            "packages": {
                "excluded": [ ],
                "included": [ ]
            },
            "withKnown": false
        },
        "sources": [ ],
        "targets": [
            "cloud-readiness"
        ]
    }
}' | jq .id)
if [ "${TASK_ID}" = "null" ]; then
  echo "Failed to create task"
  exit 1
fi
echo "Task created with id ${TASK_ID}"

# Give windup ten minutes to finish
if ! timeout 300s bash -c "until curl -S -s -X GET ${host}/tasks/${TASK_ID} | jq -e '.state == \"Succeeded\"'; do sleep 30; done"; then
  echo "##########################################"
  echo "Windup task did not complete successfully"
  echo "##########################################"
  echo "Task details"
  curl -S -s -X GET ${host}/tasks/${TASK_ID} | jq
  echo "Including pod logs"

  TASK_POD_NAMESPACED_NAME=$(curl -S -s -X GET ${host}/tasks/${TASK_ID} | jq --raw-output .pod)
  TASK_POD_NAMESPACE="${TASK_POD_NAMESPACED_NAME%/*}"
  TASK_POD_NAME="${TASK_POD_NAMESPACED_NAME#*/}"

  kubectl logs --namespace ${TASK_POD_NAMESPACE} ${TASK_POD_NAME}
  echo "Task name: ${TASK_POD_NAME} namespace: ${TASK_POD_NAMESPACE}"
  exit 0
fi
echo "Windup task completed successfully"
