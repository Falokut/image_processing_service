# Content
+ [Image Processing Service](#image-processing-service)
    + [Running instructions](#how-to-run)
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
