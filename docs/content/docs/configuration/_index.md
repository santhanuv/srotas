---
date: '2025-02-07T18:31:42+05:30'
draft: false
title: 'Configuration'
weight: 3
---

Srotas can execute a YAML configuration file to facilitate API testing and workflow automation. The configuration consists of multiple sections that define how requests are executed, variables are managed, and outputs are handled.  

{{< cards >}}
  {{< card link="/srotas/docs/configuration/global-fields" title="Global Fields" subtitle="Defines global settings such as headers and environment variables that apply across all steps." icon="globe-alt" >}}
  {{< card link="/srotas/docs/configuration/variables" title="Variables" subtitle="Covers static and dynamic variables that store and manipulate data during execution." icon="variable" >}}
  {{< card link="/srotas/docs/configuration/steps" title="Steps" subtitle="Details the different execution steps available, such as conditional logic, loops, and HTTP requests." icon="chevron-right" >}}
  {{< card link="/srotas/docs/configuration/output" title="Output" subtitle="Controls what data is returned after execution, supporting integration with chained executions." icon="external-link" >}}
{{< /cards >}}


For a detailed breakdown of each section, refer to the respective documentation.

### Sample Configuration

```yaml
base_url: "https://api.example.com/v1"

steps:
  - type: http
    step:
      name: "Fetch Users"
      method: GET
      url: "/api/users"
      store:
        users: "response.data"

  - type: forEach
    step:
      name: "Process Each User"
      list: "users"
      as: "user"
      body:
        - type: http
          step:
            name: "Fetch User Tasks"
            method: GET
            url: "/api/users/:user.id/tasks"
            store:
              tasks: "response.data"

        - type: forEach
          step:
            name: "Check Completed Tasks"
            list: "filter(tasks, {.status == 'completed'})"
            as: "task"
            body:
              - type: http
                step:
                  name: "Update Task Status"
                  method: PATCH
                  url: "/api/tasks/:task.id"
                  body:
                    status: "'reviewed'"
```

### Explanation  

#### Step 1: Fetch Users  

The first step makes a `GET` request to `/api/users` to retrieve a list of users. The response is stored in the `users` variable for later use.  

#### Step 2: Iterate Over Each User  

A `forEach` loop iterates over the `users` list. Each user is assigned to the `user` variable within the loop.  

#### Step 3: Fetch User Tasks  

For each user, another `GET` request is sent to `/api/users/:user.id/tasks` to fetch their assigned tasks. The response is stored in the `tasks` variable.  

#### Step 4: Iterate Over Completed Tasks  

Another `forEach` loop filters the `tasks` list to find those with `status == 'completed'`. Each matching task is assigned to the `task` variable.  

#### Step 5: Update Task Status  

For each completed task, a `PATCH` request is sent to `/api/tasks/:task.id` to update its status to `reviewed`.  

