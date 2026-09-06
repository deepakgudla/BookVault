output "bucket_name" {
  value = aws_s3_bucket.uploads.id
}

output "queue_url" {
  value = aws_sqs_queue.events.url
}

output "queue_arn" {
  value = aws_sqs_queue.events.arn
}

output "s3_policy_arn" {
  value = aws_iam_policy.s3_uploads.arn
}

output "sqs_policy_arn" {
  value = aws_iam_policy.sqs_events.arn
}