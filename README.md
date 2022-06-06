# cubejs-go

> The contents of this repo are a work in progress. Desired time of completion: end of April 2022.

[![Go Reference](https://pkg.go.dev/badge/github.com/TernaryInc/cubejs-go.svg)](https://pkg.go.dev/github.com/TernaryInc/cubejs-go)

## Table of Contents

- [Introduction](#introduction)
- [Support](#support)
- [Example](#example)

## Introduction

[Cube.js](https://cube.dev/) is an open-source analytics platform for data engineers and application developers to make data accessible and consistent across every application. The project includes a full-featured JS client, along with a web-based playground environment. The purpose of this client implementation in Golang is to provide an easy way to integrate from Go applications. At [Ternary](https://ternary.app/) we dogfood this client to power our users' cost-saving recommendations, anomaly detection, and billing insights.

## Support

We've intended to build the subset of functionality that we need at Ternary into the Cube client presented. We recognize that not all of the functionality available to the Javascript client has been ported, but we plan to expand feature support as find it necessary. Contributions are always welcome; see [Contributing](contributing.md).

## Future Plans

- Set up CI/CD with [GitHub Actions](https://gfgfddgleb.medium.com/how-to-test-your-go-code-with-github-actions-f15881d46089)
