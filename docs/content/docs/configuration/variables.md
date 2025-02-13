---
date: '2025-02-07T18:32:50+05:30'
draft: false
title: 'Variables'
weight: 2
---

Variables are categorized as static or dynamic based on whether other variables are available when evaluating their values. However, during execution, there is no difference between them, all variables function the same way.

### Static Variables
Static variables are initialized before execution and cannot access other variables for value evaluation. They support expressions but without variable references. The key represents the variable name, and the value is an `expr` expression that determines its initial value. These variables can be used throughout the configuration, including in global headers, and their values can be modified during execution.
```yaml
variables:
  api_key: "'secret-token-123'"
  environment: "'staging'"
  max_retries: "3"
  allowed_statuses: "['active', 'pending']"
```

| Field     | Required | Type             | Description                            |
|-----------|----------|------------------|----------------------------------------|
| variables | No       | map<string,expr> | Variables initialized before execution |

Static varibles can also be defined using the `--env` and `--var` flags. Learn more in the [Running a Configuration]({{< ref "/docs/usage/run-command.md#variables" >}}) section.

> [!CAUTION]
> Static variables are defined through the `--env` flag, `--var` flag, and the `variables` field in the configuration file. If a static variable is defined multiple times across any of these, an error is raised. 


### Dynamic Variables
Dynamic variables are created and updated as the configuration runs, allowing data to be stored and passed between steps. They can reference other variables and are typically set within different [steps]({{< ref "/docs/configuration/steps.md" >}}).

These variables are mostly created using the `store` field of an `HTTP` step and are globally available. Additionally, variables created in the `update` field of a `while` step, if they are not part of the `init` field, also persist globally, but is not recommended.

However, variables created in other step-specific fields, such as the `init` field of a `while` step or the `as` field of a `forEach` step, are only available within the scope of that step and do not persist globally.

> [!WARNING]
> Unlike static variables, which cannot be redefined, dynamic variables can be updated. If a dynamic variable already exists, its value is overwritten; otherwise, a new one is created. Static variables also become dynamic during execution.
