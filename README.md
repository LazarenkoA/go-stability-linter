# Go Dependency Stability Linter

A linter for Go that checks the stability of code dependencies in a project. The linter analyzes the Fan-in (incoming dependencies) and Fan-out (outgoing dependencies) metrics, which are used to calculate the instability (`I`) of components based on the principles outlined in Robert C. Martin's book, "Clean Architecture."

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)

## Introduction

The Go Dependency Stability Linter is a tool designed to analyze the stability of the components in your Go project. It provides insights into how your code is structured in terms of dependencies and evaluates the stability of each component based on the **Fan-in** (number of incoming dependencies) and **Fan-out** (number of outgoing dependencies) metrics.

Instability `I` is calculated using the following formula:
`I = Fan-out / (Fan-in + Fan-out)`

Where:
- `I = 0` indicates maximum stability.
- `I = 1` indicates maximum instability.

### Key Rule

The linter checks for architectural violations by enforcing the rule that **less stable components (with higher instability) must only depend on more stable components (with lower instability)**. If a more stable module depends on a less stable one, the linter will raise an error.

This ensures that your code follows a clean architecture where core components are stable and less prone to change, while higher-level components can safely depend on them without breaking stability.


## Features

- Analyzes the stability of components in Go projects.
- Provides metrics for Fan-in, Fan-out, and Instability.
- Helps identify components with high instability.
- Ensures clean architectural principles and dependency management.

## Installation

To install the Go Dependency Stability Linter, use the following command:

```bash
go install github.com/LazarenkoA/go-stability-linter@latest
```

Make sure your Go environment is set up, and `$GOPATH/bin` is added to your system's PATH.

## Usage
Once installed, you can run the linter on your Go project by executing the following command in the root of your project directory:

```bash
go-dep-stability .
```
or
```bash
go-dep-stability <progect path>
```