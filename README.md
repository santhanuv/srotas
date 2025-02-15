# Srotas

Srotas is a developer toolkit designed to simplify manual API testing by automating prerequisite steps. Instead of manually executing a series of API calls to reach your testing target, define your setup sequence in YAML. Srotas manages the execution flow, handling data between requests, so you can focus on testing your latest changes efficiently.  

## Key Features  
- **YAML-Based Configuration** – Define API requests, assertions, and workflows in a structured format.  
- **Dynamic Variables** – Store and reuse response values using expressions.  
- **Templated Requests** – Use Go `template` syntax for dynamic request bodies.  
- **Assertions & Validation** – Validate responses with powerful `expr` expressions.  
- **Debugging & Logging** – Enable `--debug` mode for detailed execution insights.  

## Getting Started  
1. Install Srotas:  

    1. Visit the [Srotas Releases](https://github.com/your-repo/srotas/releases) page.  
    2. Download the latest release for your platform (Linux, macOS, or Windows).  
    3. Move the binary to a directory in your system's `PATH`.  

    Now, you can run `srotas --help` to verify the installation.  

2. Define a test sequence in `test.yaml`:  
   ```yaml
   steps:
     - name: Fetch User Data
       request:
         method: GET
         url: "https://api.example.com/user"
       status_code: 200
   ```  
3. Run the test:  
   ```sh
   srotas run test.yaml
   ```  

## Why Srotas?  
- **Efficient Manual Testing** – The tool automates all necessary HTTP requests leading up to the feature you want to test while keeping you in control.  
- **Configurable Execution** – Customize test sequences to fit your workflow, running individual steps or complete test flows as needed.  
- **Lightweight & Portable** – No dependencies, just a single binary.

Srotas bridges the gap between manual and automated testing, giving developers a fast and structured way to test APIs without the complexity of full test automation.

Start with [**Getting Started**]({{< ref "/docs/getting-started.md" >}}) to set up Srotas quickly. Learn how to define API workflows in YAML with the [**Configuration Reference**]({{< ref "/docs/configuration" >}}), explore available commands in the [**CLI Reference**]({{< ref "/docs/usage" >}}), and check out [**Examples**]({{< ref "docs/examples.md">}}) to see practical use cases in action!

## Understanding Expressions and Templates

Srotas leverages **expr expressions** for dynamic value computation and **Go’s text/template syntax** for flexible request body definitions. Understanding these concepts will help you write powerful and customizable test configurations.

### Expr Expressions
Srotas allows you to use `expr` expressions in YAML configurations and the CLI to evaluate dynamic values at runtime. You can learn more about the syntax and capabilities of `expr` here: [Expr Language Definition](https://expr-lang.org/docs/language-definition)

### Go Text Templates
For defining HTTP request bodies, Srotas utilizes Go’s `text/template` syntax. If you're new to Go templates or want to explore advanced templating features, refer to these resources:  
[HashiCorp Nomad Go Template Guide](https://developer.hashicorp.com/nomad/tutorials/templates/go-template-syntax)  
[Gomplate Template Syntax](https://docs.gomplate.ca/syntax/)  
[Official Go text/template Documentation](https://pkg.go.dev/text/template)

