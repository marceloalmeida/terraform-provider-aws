---
subcategory: "Resource Groups Tagging"
layout: "aws"
page_title: "AWS: aws_resourcegroupstaggingapi_tag_resources"
description: |-
  Manages tags on AWS resources using the Resource Groups Tagging API.
---

# Resource: aws_resourcegroupstaggingapi_tag_resources

Manages tags on one or more AWS resources using the Resource Groups Tagging API. This resource applies the same set of tags to multiple resources.

~> **NOTE:** This resource uses the [`TagResources`](https://docs.aws.amazon.com/resourcegroupstagging/latest/APIReference/API_TagResources.html) and [`UntagResources`](https://docs.aws.amazon.com/resourcegroupstagging/latest/APIReference/API_UntagResources.html) APIs. These APIs allow you to tag resources across multiple AWS services with a single API call. However, you may prefer to use service-specific tagging resources for better type safety and validation.

## Example Usage

### Tag Multiple Resources

```terraform
resource "aws_resourcegroupstaggingapi_tag_resources" "example" {
  resource_arn_list = [
    aws_instance.example.arn,
    aws_ebs_volume.example.arn,
  ]

  tags = {
    Environment = "Production"
    Project     = "MyProject"
    ManagedBy   = "Terraform"
  }
}
```

### Tag Resources from Different Services

```terraform
resource "aws_resourcegroupstaggingapi_tag_resources" "cross_service" {
  resource_arn_list = [
    aws_s3_bucket.example.arn,
    aws_lambda_function.example.arn,
    aws_dynamodb_table.example.arn,
  ]

  tags = {
    CostCenter = "Engineering"
    Owner      = "TeamA"
  }
}
```

## Argument Reference

This resource supports the following arguments:

* `resource_arn_list` - (Required, Forces new resource) A list of ARNs of the resources that you want to apply tags to. An ARN (Amazon Resource Name) uniquely identifies a resource.
* `tags` - (Required) A map of tags to assign to the resources. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `id` - A comma-separated list of resource ARNs that this resource manages.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Resource Groups Tagging API Tag Resources using the comma-separated list of resource ARNs. For example:

```terraform
import {
  to = aws_resourcegroupstaggingapi_tag_resources.example
  id = "arn:aws:ec2:us-east-1:123456789012:instance/i-12345678,arn:aws:ec2:us-east-1:123456789012:volume/vol-12345678"
}
```

Using `terraform import`, import Resource Groups Tagging API Tag Resources using the comma-separated list of resource ARNs. For example:

```console
% terraform import aws_resourcegroupstaggingapi_tag_resources.example "arn:aws:ec2:us-east-1:123456789012:instance/i-12345678,arn:aws:ec2:us-east-1:123456789012:volume/vol-12345678"
```
