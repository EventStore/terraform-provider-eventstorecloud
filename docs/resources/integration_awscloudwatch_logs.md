---
page_title: "Resource eventstorecloud_integration_awscloudwatch_logs - terraform-provider-eventstorecloud"
subcategory: ""
description: |-
  Manages integrations for AwsCloudWatch logs.
  NOTE: This functionality is currently in beta. To access it please contact support.
---

# Resource (eventstorecloud_integration_awscloudwatch_logs)

Manages integrations for AwsCloudWatch logs.

**NOTE**: This functionality is currently in beta. To access it please contact support.

## Example Usage

```terraform
locals {
  describe_log_groups_arn = "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:*"
}

data "aws_caller_identity" "current" {}

resource "aws_cloudwatch_log_group" "esdb_logs" {
  name = "EventStoreLogs-${var.stage}"
}

resource "aws_iam_access_key" "esdb_logs" {
  user = aws_iam_user.esdb_logs.name
}

resource "aws_iam_user" "esdb_logs" {
  name = "esdb_logs_user-${var.stage}"
  path = "/esdb_logs/"
}

resource "aws_iam_user_policy" "esdb_logs" {
  name = "esdb_logs_user_policy-${var.stage}"
  user = aws_iam_user.esdb_logs.name

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:DescribeLogGroups"        
      ],
      "Effect": "Allow",
      "Resource": "${local.describe_log_groups_arn}"
    },
    {
      "Action": [        
        "logs:CreateLogStream",
        "logs:DescribeLogStreams",
        "logs:PutLogEvents"
      ],
      "Effect": "Allow",
      "Resource": "${aws_cloudwatch_log_group.esdb_logs.arn}:*"
    }    
  ]
}
EOF
}

resource "eventstorecloud_integration_awscloudwatch_logs" "cloudwatch" {
  project_id        = var.project_id
  cluster_ids       = [var.cluster_id]
  description       = "send ESDB logs to AWS CloudWatch"
  access_key_id     = aws_iam_access_key.esdb_logs.id
  secret_access_key = aws_iam_access_key.esdb_logs.secret
  group_name        = aws_cloudwatch_log_group.esdb_logs.name
  region            = var.region
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **cluster_ids** (List of String) Clusters to be used with this integration
- **description** (String) Human readable description of the integration
- **group_name** (String) Name of the CloudWatch group
- **project_id** (String) ID of the project to which the integration applies
- **region** (String) AWS region for group

### Optional

- **access_key_id** (String, Sensitive) The access key ID of IAM credentials which have permissions to create and write to the log group
- **id** (String) The ID of this resource.
- **secret_access_key** (String, Sensitive) The secret access key of IAM credentials which will be used to write to the log groups

## IAM Credentials and Security Implications

It is recommended you create credentials especially for use with this resource which have extremely limited access. A good example is shown in the snippet above, where the `aws_iam_user` resource only has permissions to describe log groups in the calling account, and can create and write streams exclusively in the log group which is also created as part of the snippet.

While is it possible to use the `eventstorecloud_integration` resource with a sink property of `awsCloudWatchLogs`, it is recommended to use the `eventstorecloud_integration_awscloudwatch_logs` resource instead as the IAM credentials get marked as sensitive to Terraform and will not be shown when running steps such as `terraform plan`. 

Even then, the IAM credentials given to this resource will be stored in the Terraform raw state as plain-text. More information on sensitive data in Terraform state can be read [here](https://www.terraform.io/language/state/sensitive-data).