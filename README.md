# AWS Appsync Builder

[![Build Status](https://travis-ci.com/ONSdigital/aws-appsync-generator.svg?branch=master)](https://travis-ci.com/ONSdigital/aws-appsync-generator) [![](https://godoc.org/github.com/ONSdigital/aws-appsync-generator/pkg/schema?status.svg)](http://godoc.org/github.com/ONSdigital/aws-appsync-generator/pkg/schema)

Generates appsync-flavour graphql schema, resolvers and terraform configuration

## Usage

| Arg             | Default          | Required | Description                    |
| --------------- | ---------------- | -------- | ------------------------------ |
| `-m --manifest` | `./manifest.yml` | no       | Manifest file to generate from |
| `-o --output`   | `./generated`    | no       | Default generated output path  |
| `-t --target`   | _none_           | yes      | `dynamo` or `sql`              |

Example:

```shell
> go run cmd/generator/main.go --m ./resources/config.yml -t dyanmo
```

## LICENSE

Copyright (c) 2019 Crown Copyright (Office for National Statistics)

Released under MIT license, see [LICENSE](LICENSE) for details.
