---
date: '2025-02-07T18:32:40+05:30'
draft: false
title: 'Global Fields'
weight: 1
---

### Version
```yaml
version: "1.0"
```

| Field   | Required | Type   | Description                                      |
|---------|----------|--------|--------------------------------------------------|
| version | Yes      | String | Configuration version identifier                |

**Description**  
The version field ensures compatibility and allows the framework to manage configuration parsing across different versions of Srotas. Currently, only "1.0" is supported.

>[!TIP]
>Always specify the version. Future versions may introduce breaking changes.

### Base URL
```yaml
base_url: "https://api.example.com/v1"
```

| Field    | Required | Type   | Description                             |
|----------|----------|--------|-----------------------------------------|
| base_url | No       | String | Base URL for all HTTP requests         |

**Description**  
The base URL provides a common prefix for all HTTP requests in the configuration. If not specified, each HTTP step must use full URLs.

>[!WARNING]
> Cannot contain dynamic expressions. Supports only static string.

> [!TIP]
> Optional, but recommended for consistent configurations.

### Global Headers
```yaml
headers:
  Accept: "'application/json'"
  Authorization: "'Bearer ' + token"
  X-Custom-Header: "'static-value'"
```

| Field   | Required | Type             | Description                            |
|---------|----------|------------------|----------------------------------------|
| headers | No       | map<string,expr> | Global headers applied to all requests |

**Description**  
Global headers are applied to every HTTP request in the configuration. Each header is defined as a key-value pair, where the key is the header name, and the value is an **expr** expression. This allows headers to be dynamically generated using static variables or values.  

If an HTTP step defines the same header as a global header, the HTTP step header takes precedence.

Global headers can also be set using the `--env` and `--header` flags. Learn more in the [Running a Configuration]({{< ref "/docs/usage/run-command.md#headers" >}}) section.

> [!WARNING]
> Use single quotes for string literals in expressions: `Accept: "'application/json'"`  
> Without single quotes, `application/json` will be interpreted as a division operation.

> [!Caution]
> Global headers are defined through the `--env` flag, `--header` flag, and the `headers` field in the configuration file. If a global header is defined multiple times across any of these, an error is raised. 


### Timeout
```yaml
timeout: 10
```


| Field Name | Type  | Required | Description |
|------------|------|----------|-------------|
| `timeout`  | int  | No       | Maximum duration (in seconds) before an HTTP request times out. Defaults to 15 seconds if not specified. |

**Description**  
The timeout field sets the maximum duration (in milliseconds) for all HTTP requests made during execution. If a request does not complete within this time, it will fail.
