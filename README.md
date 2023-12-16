# Content
+ [Image Processing Service](#image-processing-service)
    + [Running instructions](#how-to-run)
+ [Params info](#configuration-params-info)
    + [Jaeger config](#jaeger-config)
    + [Prometheus config](#prometheus-config)
+ [Documentation](#docs)
+ [Used modules](#used-modules)
+ [Author](#author)
+ [License](#license)

# Image Processing Service
The Image Processing Service is a simple image processing service with gRPC and RestAPI endpoints.

## How to Run
+ clone git repo
```shell
    git clone https://github.com/Falokut/image_processing_service.git
```
+ run by command
```shell
    docker compose up --build
```

# Configuration params info
if supported values is empty, then any type values are supported

| yml name | yml section | env name | param type| description | supported values |
|-|-|-|-|-|-|
| log_level   |      | LOG_LEVEL  |   string   |      logging level        | panic, fatal, error, warning, warn, info, debug, trace|
| healthcheck_port   |      | HEALTHCHECK_PORT  |   string   |     port for healthcheck| any valid port that is not occupied by other services. The string should not contain delimiters, only the port number|
| host   |  listen    | HOST  |   string   |  ip address or host to listen   |  |
| port   |  listen    | PORT  |   string   |  port to listen   | The string should not contain delimiters, only the port number|
| server_mode   |  listen    | SERVER_MODE  |   string   | Server listen mode, Rest API, gRPC or both | GRPC, REST, BOTH|
| enable_metrics   |      |   ENABLE_METRICS   | bool | enable metrics report or not, if true, prometheus and jaeger metrics configs will be ignored|  |
|service_name|  prometheus    | PROMETHEUS_SERVICE_NAME | string |  service name, thats will show in prometheus  ||
|server_config|  prometheus    |   | nested yml configuration  [metrics server config](#prometheus-config) | |
|jaeger|||nested yml configuration  [jaeger config](#jaeger-config)|configuration for jaeger connection ||

## Jaeger config

|yml name| env name|param type| description | supported values |
|-|-|-|-|-|
|address|JAEGER_ADDRESS|string|ip address(or host) with port of jaeger service| all valid addresses formatted like host:port or ip-address:port |
|service_name|JAEGER_SERVICE_NAME|string|service name, thats will show in jaeger in traces||
|log_spans|JAEGER_LOG_SPANS|bool|whether to enable log scans in jaeger for this service or not||

## Prometheus config
|yml name| env name|param type| description | supported values |
|-|-|-|-|-|
|host|METRIC_HOST|string|ip address or host to listen for prometheus service||
|port|METRIC_PORT|string|port to listen for  of prometheus service| any valid port that is not occupied by other services. The string should not contain delimiters, only the port number|

# Docs
+ [Swagger docs](swagger/docs/image_processing_service_v1.swagger.json)

# Used modules
+ [imaging](https://github.com/disintegration/imaging) for image encoding/decoding
+ [mimetype](https://github.com/gabriel-vasile/mimetype) for types detection
+ [bild](https://github.com/anthonynsimon/bild) for images processing

# Author

- [@Falokut](https://github.com/Falokut) - Primary author of the project

# License

This project is licensed under the terms of the [MIT License](https://opensource.org/licenses/MIT).

---
