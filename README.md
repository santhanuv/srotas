# srotas
A simple CLI for testing API.

## Features
- [x] Send HTTP requests: GET, POST, PUT, DELETE, etc.
- [x] Request chaining: Reuse values from previous responses in subsequent requests.
- [x] Configuration: Define requests using YAML files.
- [x] Headers and query params: Add custom headers and query params to HTTP requests.
- [ ] Assertions: Validate response for status, data, and headers.
- [ ] Retry and error handling: Define retry and error handling mechanisms.
- [ ] Control Structures: Define conditional, different loops.
- [ ] WebSocket: Add support for web socket.
- [ ] Config Generation: Generate config file with interactive mode.

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
        url: "/tasks"
        body:
          file: "request1.json" # location of the request JSON file
          data:
            "name": "user" # Key is the SJSON syntax for modifying request JSON body and value is the variable name
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
