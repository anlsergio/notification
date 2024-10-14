# Notification

<!-- TOC -->
* [Notification](#notification)
  * [Application Overview](#application-overview)
    * [Rate Limiting mechanism](#rate-limiting-mechanism)
    * [Idempotency](#idempotency)
  * [Development](#development)
    * [Prerequisites](#prerequisites)
    * [Getting Started](#getting-started)
    * [Helper Scripts](#helper-scripts)
    * [OpenAPI Documentation](#openapi-documentation)
    * [Testing](#testing)
<!-- TOC -->

## Application Overview

This is a Notification system that supports notifications through email.

### Rate Limiting mechanism

There's a rate-limiting mechanism in place based on a Leaky Bucket algorithm, where some rules are enforced based
on the notification type.

| Notification type | Max count | Expiration |
|-------------------|-----------|------------|
| Status            | 2         | 1 minute   |
| News              | 1         | 24 hours   |
| Marketing         | 3         | 1 hour     |

### Idempotency

This system ensures idempotency of notification message processing, meaning that no duplicates are processed in case 
of accidental message sending caused by system failure, malfunctioning of upstream sender services, or even, race conditions
in a multiple replica environment of this application itself.

The idempotency violation is based on the `correlationId` identifier that must come as part of the 
JSON request body when sending notifications through this application.

> [!IMPORTANT]
> As of now, duplicates are detected in a time span of **24 hours**, which should be enough to prevent most issues,
> meaning that, if for some reason, the same correlation ID is sent after 24 hours, **it will be considered a whole new notification**.

## Development

### Prerequisites
- Go `v1.23+`
- Make: required by the helper scripts
- Docker
- Mockery: to generate and update mocks
- Swag CLI `v1.8.4`: to generate and update OpenAPI specs*
> [!WARNING]
> Swag must be installed in a specific version of `v1.8.4` because of [some issues](https://stackoverflow.com/questions/76582283/swag-init-generates-nothing-but-general-api-information)
> recognizing annotations in dependency files.

### Getting Started

The easiest way to get started is by running the Docker stack defined in this project, so that all settings and 
dependency services are set up altogether without any overhead.

<details>

<summary>Running the application directly</summary>

First, refer to `template.env` to export the necessary environmental variables to configure the application.

Will spin up the application from your terminal
```shell
make run
```

The application will be running at `localhost:8080`
> Replace `8080` with the port you defined if you chose a different one.
```shell
curl http://localhost:8080/healthz -v
```

</details>

<details>

<summary>Running the application from a Docker container</summary>

Will spin up the application container
```shell
make docker-up
```

The application will be running at `localhost:8080`
```shell
curl http://localhost:8080/healthz -v
```

Update the docker container with your recent changes
```shell
make docker-update
```

</details>

#### Test the API

Send a notification through the API:

```shell
curl -X 'POST' \
  'http://localhost:8080/send' \
  -H 'Content-Type: application/json' \
  -d '{
  "correlationId": "0990cc56-f1b7-4f69-bc60-08fac22d41bj",
  "message": "Hey there!",
  "type": "status",
  "userId": "123-abc"
}'
```
> The correlation ID must be different for every different request.

#### Check the email inbox

Go to [MailHog's UI](http://localhost:8025/#) to check if the message has arrived to the inbox. You should be able
to see the notification message you just sent with the subject "_Status: there's a new status update_".

### Helper Scripts

For commonly used tasks and commands, there are quite a few helper commands added to the `Makefile`
of this project available for the `make` command, so make sure to check the file out to get to know the full list.

### OpenAPI Documentation

Once the application is up and running, you should be able to access the Swagger endpoint, where the OpenAPI
specifications for the routes implemented are parsed: http://localhost:8080/swagger/index.html

### Testing

Run the unit tests suite
```shell
make test
```

Run the unit tests suite while generating a coverage report
```shell
make test-cov
```

Render the test coverage report as HTML
```shell
make show-cov
```

[Back to top](#notification)
