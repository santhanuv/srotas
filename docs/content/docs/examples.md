---
date: '2025-02-08T03:10:19+05:30'
draft: false
title: 'Examples'
---

## Weather API Check

This configuration retrieves the current weather for a specified city from a weather API. The response is stored and then output as the current weather.

```yaml
version: "1.0"
base_url: "https://api.weatherapi.com"
steps:
  - type: http
    step:
      name: "Fetch Weather Data"
      method: GET
      url: "/v1/current.json"
      query_params:
        key: API_KEY
        q: city
      store:
        weather: "response.current"
      validations:
        status_code: 200
output:
  current_weather: "weather"
```

The city and API_KEY variables are provided using the --var flag, making it easy to reuse the same configuration for different cities.

{{% steps %}}

#### Step 1

The configuration sends a GET request to fetch weather data for `city`. The `city` and `API_KEY` variables are added to the `query_params` field, ensuring they are included as query parameters in the request. For example, if `API_KEY = 'abc'` and `city = 'def'`, the final URL will be `/v1/current.json?key=abc&q=def`. The response’s current weather conditions are stored in the `weather` variable and output as `current_weather`.

{{% /steps %}}

## Device Connectivity Check

This configuration pings a list of device IP addresses to check connectivity. Each device is pinged individually using a `forEach` loop.

```yaml
version: "1.0"
base_url: "http://devices.local"
variables:
  device_ips: "['192.168.1.10', '192.168.1.20', '192.168.1.30']"
steps:
  - type: forEach
    step:
      name: "Check Device Connectivity"
      list: "device_ips"
      as: "ip"
      body:
        - type: http
          step:
            name: "Ping Device"
            method: GET
            url: "/ping/:ip"
            store:
              status: "response.status"
            validations:
              status_code: 200
output:
  connectivity_status: "status"
```

{{% steps %}}

#### Step 1

A `forEach` loop iterates over the list of device IP addresses defined in `device_ips`. For each device, an HTTP GET request is sent to the `/ping/:ip` endpoint (with `:ip` replaced by the current IP). The response status is stored in `status` and later output as `connectivity_status`.

{{% /steps %}}

## Real-Time Order Status Polling

This example demonstrates real-time polling: it first initiates order processing, then uses a `while` loop to repeatedly check the order status until it is completed (or until a maximum number of attempts is reached).

```yaml
version: "1.0"
base_url: "https://api.orders.com"
variables:
  order_id: "'12345'"
steps:
  - type: http
    step:
      name: "Initiate Order Processing"
      method: POST
      url: "/orders/process"
      body:
        file: "order_request.json"
      store:
        order_status: "response.status"
      validations:
        status_code: 200

  - type: while
    step:
      name: "Monitor Order Status"
      init:
        attempts: "0"
      condition: "order_status != 'completed' && attempts < 10"
      update:
        attempts: "attempts + 1"
      body:
        - type: http
          step:
            name: "Check Order Status"
            method: GET
            url: "/orders/:order_id/status"
            store:
              order_status: "response.status"
            validations:
              status_code: 200
output:
    final_order_status: "order_status"
```

{{% steps %}}

#### Step 1

The configuration begins by initiating order processing with a POST request. It uses a JSON template from `order_request.json` and stores the initial order status in `order_status`.

#### Step 2

A `while` loop monitors the order status. The loop initializes with `attempts` set to `0`. It repeatedly sends a GET request to `/orders/:order_id/status` (with `:order_id` replaced by the value from the static variable) until the order status becomes `"completed"` or 10 attempts have been made. The final status is then output as `final_order_status`.

{{% /steps %}}

