---
date: '2025-02-07T18:33:17+05:30'
draft: false
title: 'If'
weight: 2
---

```yaml
type: if
step:
  name: "Role-based Action"
  condition: "user.role == 'admin'"
  then:
    - type: http
      step:
        name: "Admin Action"
        method: POST
        url: "/admin/action"
  else:
    - type: http
      step:
        name: "User Action"
        method: POST
        url: "/user/action"
```

| Field     | Type        | Required | Description                                    |
|-----------|-------------|----------|------------------------------------------------|
| type      | string      | Yes      | Must be `"if"`                                 |
| name      | string      | Yes      | Descriptive name for the condition             |
| condition | expr        | Yes      | Expression that evaluates to `true` or `false` |
| then      | list\<step> | Yes      | Steps to execute if `condition` is `true`      |
| else      | list\<step> | No       | Steps to execute if `condition` is `false`     |

**Description**  
Conditional steps allow branching logic based on evaluated expressions. The `condition` field is an `expr` expression that must return a boolean value. If it evaluates to `true`, the steps in `then` execute; otherwise, the steps in `else` execute (if present).

> [!NOTE]
> The `else` block is optional.

