# All declared variables. If you wan't to set custom value for any of variables,
# it's better to create file `terraform.tfvars` and put there custom configuration,
# format: `VARIABLE_NAME="variable value"`

# Access to AWS API. Create new user with `AmazonEC2FullAccess` permissions
# here `https://console.aws.amazon.com/iam/home?#/users` and put user's AWS_ACCESS_KEY_ID and
# AWS_SECRET_ACCESS_KEY to the `terraform.tfvars` file
variable "AWS_ACCESS_KEY_ID" { default = ""}
variable "AWS_SECRET_ACCESS_KEY" { default = ""}

# SSH key pair to manage and install software to EC2 instances. Don't forget to generate
# keys first: `$ ssh-keygen -f .ssh/ec2_key -N ''`
variable "KEY_PAIR_NAME"    { default = "ec2_key" }
variable "PUBLIC_KEY_PATH"  { default = "../.ssh/ec2_key.pub" }
variable "PRIVATE_KEY_PATH" { default = "../.ssh/ec2_key" }

# Settings for EC2 instances. It's better to change AWS_INSTANCE_TYPE to
# `t2.nano` if AWS Free Tier is not available for you
variable "AWS_INSTANCE_TYPE"      { default = "t2.micro" }
variable "AWS_DEFAULT_REGION"     { default = "us-west-1" }
variable "AWS_INSTANCE_AMI"       { default = "ami-8d948ced" }
variable "AWS_INSTANCE_USER_NAME" { default = "ubuntu" }
