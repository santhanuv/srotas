---
date: '2025-02-07T18:33:29+05:30'
draft: false
title: 'Foreach'
weight: 3
---

```yaml
type: forEach
step:
  name: "Batch User Processing"
  list: "filter(users, {.status == 'pending'})"
  as: "user"
  body:
    - type: http
      step:
        name: "Process User"
        method: POST
        url: "/users/process"
        body:
          file: "request.json"
          data:
            userId: "user.id"
```

| Field | Type        | Required | Description                                |
|-------|-------------|----------|--------------------------------------------|
| type  | string      | Yes      | Must be `"forEach"`                        |
| name  | string      | Yes      | Descriptive name for the iteration step    |
| list  | expr        | Yes      | Expression that evaluates to a list        |
| as    | string      | Yes      | Variable name that stores the current item |
| body  | list\<step> | Yes      | Steps to execute for each item in the list |

**Description**  
ForEach steps iterate over a dynamically generated list and execute a set of steps for each item in the list. The `list` field is an `expr` expression that must evaluate to a list. The `as` field is a string that defines the variable name used to reference the current item in each iteration.

> [!WARNING]
> The variable defined in `as` is only available within the same steps.
