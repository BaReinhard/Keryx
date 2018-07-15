# Keryx

Named after Kēryx (kay rooks), the service is an inviolatable messenger. This service provdes a standard way to interact with a chat platform. It allows for multiple platforms to be used at once or the ability to easily transition to a new platform without the end users having to change the way they interact with the service.

In more technical terminology it implements an abstraction over an adaptor pattern.

### Setup

- Setup app.yaml
- Create Bot for Slack
- Create Bot for Google Chat

```
runtime: go
api_version: go1

service: keyrx


env_variables:
 SECURE_ENDPOINT: "SECURE_ENDPOINT"
 AUTHORIZATION_HEADER: "AUTHORIZATION_HEADER"
 SLACK_TOKEN: "BOTS_KEY"
handlers:


# All URLs are handled by the Go application script
- url: /.*
  script: _go_app
```
