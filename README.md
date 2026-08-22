# ssmsh

A CLI to log in to EC2 instances via AWS Systems Manager Session Manager.

Instead of `aws ssm start-session --target <instance-id>`, it lets you log in intuitively using an EC2 Name tag.

## Install

Download a prebuilt binary from the [Releases page](https://github.com/suer/ssmsh/releases).

Or install via [mise](https://mise.jdx.dev/):

```sh
mise use -g github:suer/ssmsh
```

## Prerequisites

- [session-manager-plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html) installed and available on PATH
- AWS credentials available (`AWS_PROFILE` / `AWS_REGION`, `~/.aws/config`, etc.)

## Required IAM Policy

- `AmazonSSMManagedInstanceCore`

## Build

```sh
go build -o ssmsh .
```

## Usage

```sh
# By instance ID
ssmsh i-0123456789abcdef0

# By Name tag (resolved to an instance ID)
ssmsh web-server-01

# Run a command and exit (like ssh -c)
ssmsh web-server-01 -c "uname -a"

# Forward a local port to a remote port on the instance (like ssh -L)
ssmsh web-server-01 -L 10080:80

# Explicit profile/region
ssmsh web-server-01 --profile myprofile --region ap-northeast-1

# Print version
ssmsh -v
```

If a Name tag matches multiple instances, you can interactively pick one from the candidates.
