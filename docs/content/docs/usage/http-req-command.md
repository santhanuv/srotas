---
date: '2025-02-08T02:44:43+05:30'
draft: false
title: 'Send Http Request'
---


The `http` command in Srotas allows users to send a single HTTP request without requiring a configuration file. It is intended for quick, ad-hoc requests rather than extensive usage. This command is useful when testing an API endpoint or making a one-time request after running a configuration.

## Usage

```sh
srotas http [METHOD] [URL] [flags]
```

- **`METHOD`**: The HTTP method to use (e.g., GET, POST, PUT, DELETE).
- **`URL`**: The endpoint to send the request to.

For a complete list of available flags and usage details, run:  
```sh
srotas http --help
```

### Example
```sh
srotas http GET https://api.example.com/users
```

## Flags and Options

### Query Parameters

**Usage**  

```sh
srotas http GET https://api.example.com/users --query "key=value"
```

```sh
srotas http GET https://api.example.com/users -Q "key=value"
```

**Format**  

- `key=value`
- Multiple query parameters can be specified using multiple `--query` flags or by separating them with a comma.

**Example**   

```sh
srotas http GET https://api.example.com/users --query "status=active" --query "role=admin"
```

---

### Headers

**Usage**  

```sh
srotas http GET https://api.example.com/users --headers "key:value"
```

```sh
srotas http GET https://api.example.com/users -H "key:value"
```

**Format**  

- `key:value`
- Multiple headers can be specified using multiple `--headers` flags.

**Example**  

```sh
srotas http GET https://api.example.com/users --headers "Authorization: Bearer abc123" --headers "X-Correlation-ID: 12345"
```

---

### Request Body

**Usage**  

```sh
srotas http POST https://api.example.com/users --body '{"name": "John Doe"}'
```

```sh
srotas http POST https://api.example.com/users -B '{"name": "John Doe"}'
```

**Format**  

- A JSON-formatted string.

**Example**  

```sh
srotas http POST https://api.example.com/users --body '{"name": "Alice", "email": "alice@example.com"}' --headers "Content-Type: application/json"
```

> [!WARNING]
> The `http` command does not support expressions (`expr`). Only static values can be used.

> [!NOTE]
> It is not designed for heavy usage or complex workflows. Best used for quick testing of API endpoints or making a one-off request after running a configuration.
