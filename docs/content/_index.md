---
date: '2025-02-07T18:16:40+05:30'
draft: false
title: 'Srotas CLI'
---

Srotas is a developer toolkit designed to simplify manual API testing by automating prerequisite steps. Instead of manually executing a series of API calls to reach your testing target, define your setup sequence in YAML. Srotas manages the execution flow, handling data between requests, so you can focus on testing your latest changes efficiently.  

## Key Features  
- **YAML-Based Configuration** – Define API requests, assertions, and workflows in a structured format.  
- **Dynamic Variables** – Store and reuse response values using expressions.  
- **Templated Requests** – Use Go `text/template` for dynamic request bodies.  
- **Assertions & Validation** – Validate responses with powerful `expr` expressions.  
- **Debugging & Logging** – Enable `--debug` mode for detailed execution insights.  

## Getting Started  
1. Install Srotas:  

    1. Visit the [Srotas Releases](https://github.com/your-repo/srotas/releases) page.  
    2. Download the latest release for your platform (Linux, macOS, or Windows).  
    3. Extract the archive and move the binary to a directory in your system's `PATH`.  

    Now, you can run `srotas --help` to verify the installation.  

2. Define a test sequence in `test.yaml`:  
   ```yaml
   steps:
     - name: Fetch User Data
       request:
         method: GET
         url: "https://api.example.com/user"
       asserts:
         - "response.status_code == 200"
   ```  
3. Run the test:  
   ```sh
   srotas run test.yaml
   ```  

## Why Srotas?  
- **Efficient Manual Testing** – Automate setup steps while keeping control over the testing process.  
- **Flexible Execution** – Run individual steps or full test sequences as needed.  
- **Lightweight & Portable** – No dependencies, just a single binary.  

Srotas bridges the gap between manual and automated testing, giving developers a fast and structured way to test APIs without the complexity of full test automation.

Start with [**Getting Started**]({{< ref "/docs/getting-started.md" >}}) to set up Srotas quickly. Learn how to define API workflows in YAML with the [**Configuration Reference**]({{< ref "/docs/configuration" >}}), explore available commands in the [**CLI Reference**]({{< ref "/docs/usage" >}}), and check out [**Examples**]({{< ref "docs/examples.md">}}) to see practical use cases in action!

## Dynamic Expressions & Templating in Srotas
Srotas leverages expr expressions in its YAML configurations and CLI to compute dynamic values at runtime. For defining the HTTP request body, Srotas uses Go’s text/template syntax—see various blog posts and documentation for best practices on writing and organizing these templates.
