#!/usr/bin/env bash
 
set -euo pipefail
 
# enable debug
# set -x

echo "configuring sqs queues"
echo "======================"
LOCALSTACK_HOST=localhost
AWS_REGION=eu-central-1
 
# https://docs.aws.amazon.com/cli/latest/reference/sqs/create-queue.html
create_queue() {
  local QUEUE_NAME_TO_CREATE=$1
  awslocal --endpoint-url=http://${LOCALSTACK_HOST}:4566 sqs create-queue --queue-name ${QUEUE_NAME_TO_CREATE} --region ${AWS_REGION} --attributes VisibilityTimeout=30
}
 
create_queue "bloxroute"
