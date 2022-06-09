---
page_title: "Resource eventstorecloud_integration_awscloudwatch_logs - terraform-provider-eventstorecloud"
subcategory: ""
description: |-
  Manages AwsCloudWatch integration sink resources
---

# Resource (eventstorecloud_integration_awscloudwatch_logs)

Manages AwsCloudWatch integration sinks.

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

resource "eventstorecloud_integration_awscloudwatch_logs" "cloudwatch"{
  project_id        = var.project_id
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

- **access_key_id** (String) The access key ID of IAM credentials which have permissions to create and write to the log group
- **description** (String) Human readable description of the integration
- **group_name** (String) The name of the log group
- **project_id** (String) ID of the project to which the integration applies
- **region** (String) The AWS region name of the log group
- **secret_access_key** (String) The secret access key of IAM credentials which will be used to write to the log groups
- **source** (String) The source of the AwsCloudWatch integration sink (currently this may only be "logs")

## Import

Import is supported using the following syntax:

```shell
terraform import eventstorecloud_integration_awscloudwatch_logs.cloudwatch project_id:integration_id
```

## IAM Credentials and Security Implications

The IAM credentials given to the EventStore Cloud must ultimately be used and stored on the actual virtual machines running each cluster. For that reason it is recommended you create credentials especially for use with this resource which have extremely limited access. A good example is shown in the snippet above, where the `aws_iam_user` resource only has permissions to describe log groups in the calling account, and can create and write streams exclusively in the log group which is also created as part of the snippet.

While is it possible to use the `eventstorecloud_integration` resource with a sink property of `awsCloudWatch`, it is recommended to use the `eventstorecloud_integration_awscloudwatch_logs` resource instead as the IAM credentials get marked as sensitive to Terraform and will not be shown when running steps such as `terraform plan`. 

Even then, the IAM credentials given to this resource will be stored in the Terraform raw state as plain-text. More information on sensitive data in Terraform state can be read [here](https://www.terraform.io/language/state/sensitive-data).