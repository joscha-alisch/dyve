<p align="center">
  <img style="max-width: 960px" src="/docs/img/header.png?raw=true">
</p>

[![Go Version](https://img.shields.io/github/go-mod/go-version/joscha-alisch/dyve.svg)](https://github.com/joscha-alisch/dyve)
[![GoDoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/joscha-alisch/dyve)
[![GoReportCard](https://goreportcard.com/badge/github.com/joscha-alisch/dyve)](https://goreportcard.com/report/github.com/joscha-alisch/dyve)
[![Maintainability](https://api.codeclimate.com/v1/badges/75d5757fe6a001f6ea1b/maintainability)](https://codeclimate.com/github/joscha-alisch/dyve/maintainability)
[![Test Coverage](https://codecov.io/gh/joscha-alisch/go-cgen/branch/main/graph/badge.svg?token=898J1INMMX)](https://codecov.io/gh/joscha-alisch/dyve)
[![Sourcegraph](https://sourcegraph.com/github.com/joscha-alisch/go-cgen/-/badge.svg)](https://sourcegraph.com/joscha-alisch/dyve?badge)

# Dyve

Dyve is a vendor-agnostic unified interface for all your platforms, CI and monitoring tools.
In its basic form, Dyve allows platform operators to bundle their used tools into one solution, providing a great developer experience to their teams.
Dyve is fully extensible for more complex use-cases or tools that are not yet supported out of the box.

## Docs / How To Use

For instructions on how to install, extend and configure Dyve, please [refer to our documentation.](https://joscha-alisch.github.io/dyve)

## Demo

To evaluate whether Dyve is a fit for your company, you can run docker-compose

## Supported Tools and Platforms

The following products are supported in Dyve. Checked items are fully implemented and usable while the others are currently in development.
If you want to help out, refer to the [corresponding issues in GitHub]() or [create a new one]() if you find something missing here.

### Platforms

* [x] CloudFoundry
* [ ] Kubernetes
* [ ] Google Cloud - Compute Engine

### Continuous Integration & Deployment 

* [x] Concourse
* [ ] GitHub Actions

### Monitoring
* [ ] Sentry
* [ ] Prometheus
* [ ] ElasticSearch / Kibana
* [ ] Google Cloud Monitoring

