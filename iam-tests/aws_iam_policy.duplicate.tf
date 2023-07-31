resource "aws_iam_policy" "duplicate_keys" {
  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": [
      "iam:Get*",
      "iam:List*"
    ],
    "Action": [
      "iam:Delete*"
    ],
    "Resource": "*"
  }
}
POLICY
}