---
date: '2025-02-07T18:33:25+05:30'
draft: false
title: 'While'
weight: 4
---

```yaml
type: while
step:
  name: "Job Status Monitoring"
  init:
    attempts: "0"
    status: "'pending'"
  condition: "status != 'completed' && attempts < 10"
  update:
    attempts: "attempts + 1"
  body:
    - type: http
      step:
        name: "Check Status"
        method: GET
        url: "/jobs/:job_id"
        store:
          status: "response.status"
```

| Field     | Type             | Required | Description                                           |
|-----------|------------------|----------|-------------------------------------------------------|
| type      | string           | Yes      | Must be `"while"`                                     |
| name      | string           | Yes      | Descriptive name for the polling step                 |
| init      | map<string,expr> | No       | Variables initialized before the loop starts          |
| condition | expr             | Yes      | Expression that determines whether the loop continues |
| update    | map<string,expr> | No       | Variables updated after each iteration                |
| body      | list\<step>      | Yes      | Steps to execute while the condition is true          |

**Description**  
While steps repeatedly execute a set of steps until a condition evaluates to false. The `init` field is a map where the key is the variable name and the value is an `expr` expression. These variables are initialized before the loop starts and are only available within the same while step.

The `condition` field is an `expr` expression that determines whether the loop should continue execution and must evaluate to a boolean value. It is evaluated after each iteration.

The `update` field is also a map where the key is the variable name and the value is an `expr` expression. It can reference all available variables, including the ones defined in `init`. The `update` expressions is evaluated after each iteration.

> [!WARNING]
> The `init` field defines variables that exist only within the `while` step. Any modifications to `init` variables are lost outside the `while` step since they are scoped only to that step.

> [!TIP] 
> To ensure values does not persist beyond the `while` step, only `init` variables should be used within `update` instead of creating new ones unintentionally. Variables created in `update` persist globally, but modifications to `init` variables do not.  
