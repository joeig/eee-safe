# eee-safe

A custom Threema Safe server implementation with multi-backend support.

[![Build Status](https://travis-ci.org/joeig/eee-safe.svg?branch=master)](https://travis-ci.org/joeig/eee-safe)
[![Go Report Card](https://goreportcard.com/badge/github.com/joeig/eee-safe)](https://goreportcard.com/report/github.com/joeig/eee-safe)

Currently supported backends:

* DynamoDB
* Filesystem

## Client setup

Threema supports third party backup servers out of the box. Simply choose "Threema Safe", "Expert settings" and deactivate "Use default server". You will be asked for your custom server endpoint and credentials.

## Server setup

### Install from source

You need `go` and `GOBIN` in your `PATH`. Once that is done, install `eee-safe` using the following command:

~~~ bash
go get -u github.com/joeig/eee-safe
~~~

### Configuration

After that, copy [`config.dist.yml`](configs/config.dist.yml) to `config.yml`, replace the default settings and run the binary:

~~~ bash
eee-safe -config=/path/to/config.yml
~~~

If you're intending to add the application to your systemd runlevel, you may want to take a look at [`init/eee-safe.service`](init/eee-safe.service).

Threema requires a valid CA certificate.

Choose one of the following storage backends:

Storage backend | DynamoDB | Filesystem
--------------- | -------- | ----------
Built-in TTL    | yes      | no
Thread safe     | yes      | no
Works without further efforts | no | yes

#### DynamoDB settings

This option requires a pre-configured AWS environment.

##### Table settings

Create a new table with the following settings:

* Table name in this example: `EEESafe`
* Primary key: `backupID`
* Time to live attribute: `expirationTime`
* Turn backup functionality on if required.

##### IAM

Required IAM permissions to access the DynamoDB table (change the resource ARN if necessary):

~~~ json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "MaintainEEESafeTable",
            "Effect": "Allow",
            "Action": [
                "dynamodb:PutItem",
                "dynamodb:GetItem",
                "dynamodb:DeleteItem"
            ],
            "Resource": "arn:aws:dynamodb:eu-west-1:1234567890:table/EEESafe"
        }
    ]
}
~~~

Examples:

- Provide credential files: `~/.aws/credentials` and `~/.aws/config`
- Set credentials via environment variables:
  ~~~ bash
  export AWS_REGION="eu-west-1"
  export AWS_ACCESS_KEY_ID="Your Access Key ID"
  export AWS_SECRET_ACCESS_KEY="Your Secret Access Key"
  ~~~
- Use EC2 instance roles/ECS task roles/Lambda roles/... (you should always choose this option whenever possible!)

See also: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html

#### Filesystem settings

This option stores every backup in a dedicated file on the local filesystem.

This storage backend does currently **not** support thread-safety and auto-deletion of expired backups. You probably want to implement auto-deletion by using a `find` cronjob.

## Troubleshooting

Run `eee-safe` in debug mode in order to increase the verbosity tremendously: `-debug`

## Contribution

Feel free to contribute. This is the API reference: [Cryptography Whitepaper](https://threema.ch/press-files/2_documentation/cryptography_whitepaper.pdf)

This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout) principals.