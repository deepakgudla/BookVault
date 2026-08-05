#!/bin/bash

#bucket
awslocal s3 mb s3://bookvault-uploads

#AWS SQS Queue
awslocal sqs create-queue --queue-name bookvault-events

echo "localstack initialization complete"
