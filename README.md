# AWS Appsync Builder

![Unit tests](https://github.com/ONSdigital/aws-appsync-generator/workflows/Unit%20tests/badge.svg) [![](https://godoc.org/github.com/ONSdigital/aws-appsync-generator/pkg/schema?status.svg)](http://godoc.org/github.com/ONSdigital/aws-appsync-generator/pkg/schema) [![Go Report Card](https://goreportcard.com/badge/github.com/ONSdigital/aws-appsync-generator)](https://goreportcard.com/report/github.com/ONSdigital/aws-appsync-generator)

Generates appsync-flavour graphql schema, resolvers and terraform configuration

## Usage

| Arg             | Default          | Required | Description                                                                                               |
| --------------- | ---------------- | -------- | --------------------------------------------------------------------------------------------------------- |
| `-m --manifest` | `./manifest.yml` | no       | Manifest file to generate from                                                                            |
| `-o --output`   | `./generated`    | no       | Default generated output path **Warning: Anything existing in this path will be wiped before generation** |

Example:

```shell
> go run cmd/generator/main.go --m ./resources/config.yml
```

## Manifest reference

The manifest schema reference can be found in the [documentation folder](docs/manifest-reference.md)

## LICENSE

Copyright (c) 2019 Crown Copyright (Office for National Statistics)

Released under MIT license, see [LICENSE](LICENSE) for details.
