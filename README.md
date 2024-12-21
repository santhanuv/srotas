# srotas
A developer toolkit designed to streamline API endpoint testing by automating prerequisite steps. Instead of manually executing a series of API calls to reach your testing target, define your setup sequence in YAML. Srotas handles the API chain execution, managing data flow between requests, allowing you to focus on testing your newly implemented features. Perfect for developers who need efficient manual testing workflows without the overhead of full test automation.

## Features
- [x] Send HTTP requests: GET, POST, PUT, DELETE, etc.
- [x] Request chaining: Reuse values from previous responses in subsequent requests.
- [x] Configuration: Define requests using YAML files.
- [x] Headers and query params: Add custom headers and query params to HTTP requests.
- [x] Assertions: Validate response for status and response data.
- [ ] Control Structures: Define conditional, different loops.
- [ ] WebSocket: Add support for web socket.
- [ ] Config Generation: Generate a config file in interactive mode.

## YAML configuration
```yaml
version: "1.0"
base_url: "<url>" # Base URL for all relative endpoints
timeout: 8000 # Request timeout in milliseconds
headers: # Common headers for all requests.
  Content-Type: "application/json"

sequence:
  name: "Happy Path"
  description: "Description"
  steps:
    - type: http
      step:
        name: "Get user details"
        method: GET
        url: "/users" # If base_url is specified, the URL defines the endpoint; otherwise, a full URL can be provided, which takes precedence
        delay: 3000
        store: # Extracts and saves data from the response
          user: "2.name" # Key is the name of the variable and value is the GJSON syntax for extracting JSON response value.
    - type: http
      step:
        name: "Create a task"
        method: POST
        headers:
          Content-Type: "application/json"
        url: "/:user/tasks" # URL parameters can substitute variable values. The parameter name and variable name should match.
        body:
          file: "request1.json" # location of the request JSON file
          data:
            "name": "user" # Key is the SJSON syntax for modifying request JSON body and value is the variable name
        validations: # Validate response
          status_code: 201 # Validate the status code of the response
          asserts:
            - value: "Task created successfully"
              selector: "message" # Checks if the response message matches the expected value
            - value: "$user"
              selector: "created_by" # Validates the 'created_by' field matches the 'user' variable (USE GJSON syntax)
```

[GJSON](https://github.com/tidwall/gjson) is used to extract data from response JSON. Store property expects the key values to be in GJSON syntax. [SJSON](https://github.com/tidwall/sjson) is used to add or replace JSON request values with variable values from the store. The data property of Body expects the key to be of SJSON syntax.

## Usage
Build the binary using go cli.

### Run configuration
```bash
$ srotas run CONFIGURATION.yaml
```

### Run HTTP request
```bash
$ srotas http GET http://localhost:8080/task
```
See `--help` to know more about the commands.
