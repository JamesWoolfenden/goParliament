resource "aws_iam_policy" "duplicate_keys_2" {
  policy = data.aws_iam_policy_document.duplicate.json
}

data "aws_iam_policy_document" "duplicate" {
  statement {
    sid="1"
  
    actions = [
      "s3:ListAllMyBuckets",
      "s3:GetBucketLocation",
    ]

    resources = [
      "arn:aws:s3:::*"
    ]

    actions = [
      "ec2:*"
    ]
  }
}