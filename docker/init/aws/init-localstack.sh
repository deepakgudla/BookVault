#!/bin/bash

awslocal s3 mb s3://bookvault-uploads
awslocal s3api put-bucket-versioning \
	--bucket bookvault-uploads \
	--versioning-configuration Status=Enabled
awslocal s3api put-bucket-encryption \
	--bucket bookvault-uploads \
	--server-side-encryption-configuration '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'

awslocal sqs create-queue --queue-name bookvault-events-dlq
DLQ_ARN=$(awslocal sqs get-queue-attributes \
	--queue-url http://localhost:4566/000000000000/bookvault-events-dlq \
	--attribute-names QueueArn \
	--query 'Attributes.QueueArn' \
	--output text)
awslocal sqs create-queue \
	--queue-name bookvault-events \
	--attributes "VisibilityTimeout=60,RedrivePolicy={\"deadLetterTargetArn\":\"${DLQ_ARN}\",\"maxReceiveCount\":\"5\"}"

echo "localstack initialization complete"
