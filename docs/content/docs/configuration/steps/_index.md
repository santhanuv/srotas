---
date: '2025-02-07T18:33:09+05:30'
draft: false
title: 'Steps'
weight: 3
---

Steps define the execution flow within a Srotas configuration. Each step represents an action or decision point, such as making an HTTP request, iterating over a list, or executing conditional logic.  

A step consists of the following fields:  

| Field       | Type   | Required | Description                                                     |
|-------------|--------|----------|-----------------------------------------------------------------|
| `type`      | string | Yes      | Defines the step type (e.g., `http`, `if`, `while`, `forEach`). |
| `step`      | object | Yes      | Contains step-specific configurations.                          |
| `step.name` | string | Yes      | A descriptive name for the step.                                |

### Step Types

Srotas supports multiple step types, each serving a specific purpose:  

{{< cards >}}
  {{< card link="/docs/configuration/steps/http" title="HTTP Step" icon="server" >}}
  {{< card link="/docs/configuration/steps/if" title="If Step" icon="filter" >}}
  {{< card link="/docs/configuration/steps/foreach" title="ForEach step" icon="duplicate" >}}
  {{< card link="/docs/configuration/steps/while" title="While Step" icon="refresh" >}}
{{< /cards >}}

Each step type has its own structure and fields, detailed in their respective sections.  


> [!NOTE]
> Steps execute sequentially unless modified by control structures like `if`, `forEach`, or `while`.  

> [!NOTE]
> Variables created in a step may be available to subsequent steps, depending on scope rules.  
