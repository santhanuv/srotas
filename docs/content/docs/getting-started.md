---
date: '2025-02-07T18:34:14+05:30'
draft: false
title: 'Getting Started'
weight: 1
---

Srotas is a command-line tool for testing APIs using YAML configuration files. This guide will help you get started quickly.

## Installation

To install Srotas, visit the [GitHub Releases](https://github.com/santhanuv/srotas/releases) page and download the latest release for your platform. Follow the installation instructions provided.

## Running a Configuration

Once installed, you can run a Srotas configuration file using:

```sh
srotas run config.yaml
```

### Common Flags

- `-D, --debug` : Enables debug mode for detailed logs.
- `-E, --env <file or JSON>` : Loads global headers and variables.
- `-H, --header <key:value>` : Adds an additional global header.
- `-V, --var <name=value>` : Defines a global variable.

### Example Usage

Run a configuration with debug mode:

```sh
srotas run config.yaml -D
```

Load variables and headers from a JSON file:

```sh
srotas run config.yaml -E env.json
```

Pass headers dynamically:

```sh
srotas run config.yaml -H "Authorization: 'Bearer ' + token"
```

Define variables inline:

```sh
srotas run config.yaml -V "username=example"
```

### Chaining Configurations with Piping

Srotas allows piping outputs between executions:

```sh
srotas run fetch_users.yaml | srotas run process_tasks.yaml
```

This allows the output of the first configuration to be used as static variables in the second.

For more advanced usage, check out the [Detailed Usage Guide]({{< ref "/docs/usage.md" >}}).

