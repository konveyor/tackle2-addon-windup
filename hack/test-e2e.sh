#!/bin/bash

host="${HOST:-localhost:8080}"

echo "$(kubectl wait deployment/tackle-hub --for condition=available --timeout=-1s -n konveyor-tackle)"

# Create a Stake Holder Group
curl -X POST ${host}/stakeholdergroups -d \
'{
    "name": "Big Dogs",
    "description": "Group of big dogs."
}' | jq -M .

# Create a Stake Holder
curl -X POST ${host}/stakeholders -d \
'{
    "name": "tackle",
    "displayName":"Elmer",
    "email": "tackle@konveyor.org",
    "role": "Administrator",
    "stakeholderGroups": [{"id": 1}],
    "jobFunction" : {"id": 1}
}' | jq -M .

# Create a Business Service
curl -X POST ${host}/businessservices -d \
'{
    "createUser": "tackle",
    "name": "Marketing",
    "Description": "Marketing Dept.",
    "owner": {
      "id": 1
    }
}' | jq -M .

# Create an Application
curl -X POST ${host}/applications -d \
'{
    "createUser": "tackle",
    "name":"Pathfinder",
    "description": "Tackle Pathfinder application.",
    "repository": {
      "name": "tackle-pathfinder",
      "url": "https://github.com/konveyor/tackle-pathfinder.git",
      "branch": "1.2.0"
    },
    "facts": {
      "analysed": true
    },
    "businessService": {"id":1}
}' | jq -M .

# Create a Review
curl -X POST ${host}/reviews -d \
'{
    "businessCriticality": 4,
    "effortEstimate": "extra_large",
    "proposedAction": "repurchase",
    "workPriority": 1,
    "comments": "This is hard.",
    "application": {"id":1}
}' | jq -M .

# Make a request to hub
request_cmd="$(curl -i -o - -X POST ${host}/tasks -d \
'{
    "name":"Windup",
    "state": "Ready",
    "locator": "windup",
    "addon": "windup",
    "application": {"id": 1},
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
}')"

# Get status code from the curl request
status_code="$(echo "$request_cmd" | grep HTTP | awk '{print $2}')"

# Get output from the curl request
output_response=$(echo "$request_cmd")
echo "Output response: $output_response"

# Check if status_code starts with 2
if [[ "${status_code}" != 2* ]]; then
    echo "Failed to create windup task"
    echo "Got Response Status: ${status_code}"
    exit 1
fi

pods="$(kubectl get pods -n konveyor-tackle)"
task_status="$(echo "$pods" | grep task | awk '{print $3}')"
if [[ $task_status != "Completed" ]]; then
  echo "Failed to create windup task"
  echo "Got Task Pod Status: $task_status"
  exit 1
fi
