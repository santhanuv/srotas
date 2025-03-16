---
date: '2025-02-07T18:33:34+05:30'
draft: false
title: 'HTTP Request'
weight: 1
---

An HTTP step defines an API request with options for dynamic parameters, request body templating, and response validation.  

```yaml
type: http
step:
  name: "Fetch User Details"
  method: GET
  url: "/users/:user_id"
  query_params:
    key: "value"
  headers:
    Custom-Header: "dynamic_value"
  body:
    file: "request_template.json"
    data:
      dynamic_field: "computed_value"
  store:
    user_data: "response.user"
  validations:
    status_code: 200
    asserts:
      - "response.status == 'online'"

```

| Field                   | Type                      | Required | Description                                                     |
|-------------------------|---------------------------|----------|-----------------------------------------------------------------|
| type                    | string                    | Yes      | Must be `"http"`                                                |
| name                    | string                    | Yes      | Descriptive name for the step                                   |
| method                  | string                    | Yes      | HTTP method (`GET`, `POST`, etc.)                               |
| url                     | string                    | Yes      | Request URL with optional path parameters                       |
| headers                 | map<string, list\<expr\>> | No       | Request-specific headers                                        |
| query_params            | map<string, list\<expr\>> | No       | URL query parameters                                            |
| delay                   | int                       | No       | Delays the HTTP request execution by the specified time (in ms) |
| body.file               | string                    | No       | JSON template file path                                         |
| body.template           | string                    | No       | Inline JSON template                                            |
| body.data               | map<string, expr>         | No       | Dynamic data for template                                       |
| store                   | map<string, expr>         | No       | Variables to extract from the response                          |
| validations.status_code | int                       | No       | Expected HTTP status code                                       |
| validations.asserts     | list\<expr>               | No       | List of validation expressions                                  |

#### URL Parameters

Use `/:variable` in the URL to specify path parameters. These parameters are replaced with values from stored variables.

For example, if userId is a variable with the value 123, the following URL:

```yaml
url: "/users/:userId"
```

Will be transformed into: `/users/123`

> [!WARNING]
> `expr` expressions are not allowed in URL parameters.

> [!CAUTION]
> The variable used in the parameter must already be defined; otherwise, an error will be raised.

#### Headers

Defined as a map where the key is the header name and the value is a list of `expr` expression that evaluate to a string. All variables can be used in these expressions. Step headers are prioritized over global headers if a header with the same name exists.

#### Body

The request body can be defined using either a file or an inline template. If both are provided, the template takes precedence.

##### Body Fields

- `data`:  
Specifies the input for the template. Supports expr expressions for dynamic computation. All defined data fields are accessible within the template.

- `file`:  
Path to an external file containing the request body template.

- `template`:  
Inline request body template. If provided, this is used instead of the file.

> [!IMPORTANT]
> The template for the `body` field should follow the Go text/template format.

#### Store

A map where the key is the variable name, and the value is an `expr` expression. The `response` variable holds the HTTP response body and can be used only within the same HTTP step.

#### Validations

- `status_code` specifies the expected HTTP status code.  
- `asserts` is a list of `expr` expressions that validate the response body. These expressions have access to all variables and the `response` variable and must return a boolean value.

> [!NOTE]
> `store` captures response data after the request completes.  

### HTTP Request Template

When it comes to defining the HTTP request body, Srotas uses Go’s built-in `text/template` syntax. This lets you create a template that mixes static JSON with dynamic data. You can define inline templates or reference external files, and you have full control over how the final JSON is generated. When specifying the request template in a file, ensure that the main template is defined using {{define "request"}} ... {{end}}. This template is used as the HTTP request body.

**Learning Go Template Syntax**

If you're new to Go templates, follow these resources in order to get familiar with the syntax:  

1. [Nomad’s Go Template Syntax Guide](https://developer.hashicorp.com/nomad/tutorials/templates/go-template-syntax) – A beginner-friendly introduction with examples.  
2. [Gomplate Syntax Documentation](https://docs.gomplate.ca/syntax/) – Covers basic templating.  
3. [Official Go text/template Documentation](https://pkg.go.dev/text/template) – The full reference for Go's templating engine.  

By following these, you'll quickly get up to speed with writing templates for Srotas.  
