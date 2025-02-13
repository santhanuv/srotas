---
date: '2025-02-07T18:32:54+05:30'
draft: false
title: 'Output'
weight: 4
---

```yaml
output:
  processed_items: "length(items)"
  success_rate: "processed_count / total_count * 100"
  final_status: "last_response.status"
```

| Field      | Required | Type             | Description                                                        |
|------------|----------|------------------|--------------------------------------------------------------------|
| output     | No       | map<string,expr> | Defines variables to be included in the final output               |
| output_all | No       | bool             | If true, includes all variables in output, ignoring `output` field |

**Description**  
The `output` section defines the final result of the configuration execution. It allows selecting specific variables to include in the output. The `output` field is a map where the key is the output name and the value is an `expr` expression that is evaluated at the end of execution. All variables available during execution can be used in these expressions.  

Additionally, the `output_all` field can be used to include all available variables in the output. If `output_all` is set to `true`, the `output` field is ignored, and the entire execution context is added to the output data.  

The output data is particularly useful when using the pipe (`|`) operator to chain multiple Srotas configurations. When output is piped, the next configuration execution reads it as static variables, effectively transferring the execution context between runs.

> [!NOTE]
> Output is generated only after the entire configuration has completed execution.

> [!WARNING]
> If both `output` and `output_all` are provided, `output_all` takes precedence.

