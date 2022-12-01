#!/bin/bash

set -o errexit
set -o nounset

host="${HOST:-localhost:8080/hub}"

# Port Forwarding
kubectl port-forward service/tackle-ui 8080:8080 -n konveyor-tackle > /dev/null 2>&1 &
pid=$!

# kill the port-forward regardless of how this script exits
trap '{
    kill $pid
}' EXIT

# wait for port to become available
while ! nc -vz localhost 8080 > /dev/null 2>&1 ; do
    sleep 0.1
done

# Create a Stake Holder Group if one does not exist
stakeholder_groups="$(curl -X GET ${host}/stakeholdergroups)"
if [[ "${stakeholder_groups}" == "[]" ]]; then
    stakeholder_group="$(curl -X POST ${host}/stakeholdergroups -d \
        '{
            "name": "Big Dogs",
            "description": "Group of big dogs."
        }')"
    echo "Stakeholder Group ID: $(echo "$stakeholder_group" | jq -r '.id')"
fi

# Create a Stake Holder if one does not exist
stakeholders="$(curl -X GET ${host}/stakeholders)"
if [[ "${stakeholders}" == "[]" ]]; then
    stakeholder="$(curl -X POST ${host}/stakeholders -d \
        '{
            "name": "tackle",
            "displayName":"Elmer",
            "email": "tackle@konveyor.org",
            "role": "Administrator",
            "stakeholderGroups": [{"id": 1}],
            "jobFunction" : {"id": 1}
        }')"
    echo "Stakeholder ID: $(echo "$stakeholder" | jq -r '.id')"
fi

# Create a Business Service if one does not exist
business_services="$(curl -X GET ${host}/businessservices)"
if [[ "${business_services}" == "[]" ]]; then
    business_service="$(curl -X POST ${host}/businessservices -d \
        '{
            "createUser": "tackle",
            "name": "Marketing",
            "Description": "Marketing Dept.",
            "owner": { "id": 1}
        }')"
    echo "Business Service ID: $(echo "$business_service" | jq -r '.id')"
fi

# Create an Application
application_data='{
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
}'

application="$(curl -X POST ${host}/applications -d "${application_data}")"
application_id="$(echo "$application" | jq -r '.id')"
# If application already exists, find the application by name in the list of applications and get the ID
if [[ "${application_id}" == "null" ]]; then
    application_name="$(echo "$application_data" | jq -r '.name')"
    application_id="$(curl -X GET ${host}/applications | jq -r --arg application_name "${application_name}" '.[] | select(.name == $application_name) | .id')"
fi
echo "Application ID: $application_id"

# Make a request to hub
request_cmd="$(curl -iX POST ${host}/tasks -d \
'{
    "name":"Windup",
    "state": "Ready",
    "locator": "windup",
    "addon": "windup",
    "application": {"id": '$application_id'},
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
output_response=$(echo "$request_cmd" | awk '/^{/{p=1} p')
echo "$output_response"

# Check if status_code starts with 2
if [[ "${status_code}" != 2* ]]; then
    echo "Failed to create windup task"
    echo "Got Response Status: ${status_code}"
    exit 1
fi

# Get task id from the output response
task_id="$(echo "$output_response" | jq -r '.id')"

# Get the task pod with task_id
sleep 3
task_pod="$(kubectl get pods -n konveyor-tackle | grep task-$task_id | awk '{print $1}')"
echo "waiting for $task_pod to be completed..."

# Wait for the task pod to succeed or fail
while true; do
    task_status="$(kubectl get pods -n konveyor-tackle $task_pod -o jsonpath='{.status.phase}')"
    if [[ "${task_status}" == "Succeeded" ]]; then
        echo "Task Succeeded"
        break
    elif [[ "${task_status}" == "Failed" ]]; then
        echo "Failed to complete windup task"
        exit 1
    fi
    task_pod="$(kubectl get pods -n konveyor-tackle | grep task-$task_id | awk '{print $1}')"
    sleep 5
done
