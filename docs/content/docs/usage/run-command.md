---
date: '2025-02-07T18:34:14+05:30'
draft: false
title: 'Running a Configuration'
weight: 1
---

## Usage
Use the `run` command to execute a configuration file:
```sh
srotas run [CONFIG]
```
`config`: The path to the yaml configuration file.

For more details on configuring Srotas, check out the [Configuration]({{< ref "/docs/configuration" >}}).

For a complete list of available flags and usage details, run:  
```sh
srotas run --help
```

### Example
```sh
srotas run config.yaml
```

## Flags and Options

### Debug Mode
Enable detailed logs:
```sh
srotas run config.yaml --debug
```
Alias:
```sh
srotas run config.yaml -D
```

### Environment

The `--env` flag is used to load global headers and static variables from a JSON string or a file containing JSON data.

**Usage**  
```sh
srotas run config.yaml --env env.json
```
Alias:
```sh
srotas run config.yaml -E env.json
```

**Format**  

The JSON structure consists of two fields:  

- **`Variables`**: A map where keys are variable names and values are expressions. These variables are static, meaning they do not have access to other values at the time of definition.
- **`Headers`**: A map where keys are header names and values are lists of expressions. The expressions for headers can reference static variables.

You can provide the JSON data directly as a string or specify a file path containing the JSON.

```json
{
  "Variables": {
    "auth_token": "'Bearer abc123'"
  },
  "Headers": {
    "Authorization": ["auth_token"],
  }
}
```

**Example**  

```sh
srotas run --env '{"Variables":{"user":"John Doe","token":"ey1234"},"Headers":{"Authorization":["\"Bearer \" + token"]}}' config.yaml
```

```sh
srotas run --env env.json config.yaml
```

> [!IMPORTANT]
> At least one of `Variables` or `Headers` must be present in the JSON.

> [!WARNING]
> Duplicate variable and header names will result in an error.

For more details, refer [Global Fields]({{< ref "/docs/configuration/global-fields.md" >}}) and [Variables]({{< ref "/docs/configuration/variables.md" >}}).

### Headers

The `--header` flag is used to define global headers dynamically when running a Srotas configuration. Headers specified using this flag override any global headers defined in the configuration file.

**Usage**  

```sh
srotas run --header "key1:'value1'" --header "key2: 'value2'" config.yaml
```
Alias:
```sh
srotas run -H "key1:'value1'" -H "key2: 'value2'" config.yaml
```

**Format**  

The `--header` flag follows the format:

```
"key:value"
```

- **`key`**: The name of the header.
- **`value`**: An expression that evaluates to a string. It has access to all static variables.

Multiple headers can be specified by using the flag multiple times.

**Examples**  

Define an `Authorization` header dynamically:

```sh
srotas run --header "Authorization: 'Bearer $TOKEN'" config.yaml
```

> [!TIP]
> This flag is useful when adding authentication headers dynamically instead of modifying the configuration file.

> [!WARNING]
> Multiple headers with the same name **cannot** be defined globally (in `--header`, `--env`, or the config file). If a duplicate header is defined, an error is raised.

> [!NOTE]
> Headers specified using `--header` take precedence over those defined in `--env` and the configuration file.

For more details, refer [Global Headers]({{< ref "/docs/configuration/global-fields.md#global-headers" >}}).



### Variables

The `--var` flag in Srotas allows users to define global variables directly through the command line. It is useful for setting static values that can be referenced within the configuration file.

**Usage**   

```sh
srotas run --var key1=value1 --var key2=value2 config.yaml
```
Alias:
```sh
srotas run -V key1=value1 -V key2=value2 config.yaml
```

**Format**  
Each `--var` flag must follow the format:

```
key=value
```

- **key**: The name of the variable. Must be unique.
- **value**: An [expr](https://expr-lang.org/docs/language-definition) expression. This expression is evaluated without access to any other variables.

Multiple variables can be specified by using the flag multiple times.

**Examples**  

```sh
srotas run --var timeout=30 --var api_version="v2" config.yaml
```

> [!NOTE]
> **No variable references**: The expressions in `--var` do not have access to any static or dynamic variables.

> [!WARNING]
> If a duplicate variable name is detected, Srotas will return an error.

For more details, refer [Variables]({{< ref "/docs/configuration/variables.md#static-variables" >}}).


## Chaining Configurations
Srotas supports piping output between executions:

```sh
srotas run fetch_users.yaml | srotas run process_tasks.yaml
```

This transfers the execution context from the first configuration to the second. Only the specified [Output]({{< ref "/docs/configuration/output.md" >}}) will be transfered.

> [!NOTE]
> The `output` or `output_all` field is mandatory for proper configuration chaining.
