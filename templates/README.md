# Ermes-labs Classic templates

[![Build Status](https://github.com/ermes-labs/templates/workflows/ci-only/badge.svg?branch=master)](https://github.com/ermes-labs/templates/actions)

To find out more about the OpenFaaS templates see the [faas-cli](https://github.com/openfaas/faas-cli).

> Note: The templates are completely customizable - so if you want to alter them please do fork them and use `ermes-cli template pull https://github.com/ermes-labs/templates/` to make use of your updated versions.

### Classic Templates

This repository contains the Classic Ermes-labs templates, but many more are available in the Template Store. Read above for more information.

| Name           | Language | Version | Linux base   | Watchdog | Link                                                                                       |
| :------------- | :------- | :------ | :----------- | :------- | :----------------------------------------------------------------------------------------- |
| ermes-go       | Go       | 1.22    | Alpine Linux | classic  | [Go template](https://github.com/ermes-labs/templates/tree/master/template/ermes-go)       |
| ermes-go-redis | Go       | 1.22    | Alpine Linux | classic  | [Go template](https://github.com/ermes-labs/templates/tree/master/template/ermes-go-redis) |

For more information on the templates check out the [docs](https://docs.openfaas.com/cli/templates/).

### Classic vs of-watchdog templates

The current version of OpenFaaS templates use the original `watchdog` which `forks` processes - a bit like CGI. The newer watchdog [of-watchdog](https://github.com/openfaas-incubator/of-watchdog) is more similar to fastCGI/HTTP and should be used for any benchmarking or performance testing along with one of the newer templates. Contact the project for more information.

### Contribute to this repository

See [contributing guide](https://github.com/ermes-labs/templates/blob/master/CONTRIBUTING.md).

### License

This project is part of the Ermes-labs project licensed under the MIT License.
