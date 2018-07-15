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

### Requirements

- Authorization Bearer Header
- Space Header
  - Space for Slack (slug after /messages/): ![Slug after /messages/](https://user-images.githubusercontent.com/13072194/42738885-84bb6ff2-8840-11e8-9d05-e97c4bf798ff.png)
  - Space for Google Chat (slug after /room/): ![Slug after /room/](https://user-images.githubusercontent.com/13072194/42738889-93a9ad1c-8840-11e8-9a25-4c809b745632.png)
- Destination Header
- Payload Body

### Optional (For Google Chat Only)

- Thread Header
- ThreadKey Header

### Google Example

```
curl -X POST https://YOUR_GOOGLE_URL.appspot.com/SECURE_ENDPOINT -H "Authorization: Bearer $TOKEN" \
-H "Destination: google" \
-H "Space: AAAAsDtCfAE" \
-H "ThreadKey: PR-2" \
-H "Content-Type: application/json" \
-d '{"text":"PR Opened on Blah Blah"}'
```

### Slack Example

```
curl -X POST https://YOUR_GOOGLE_URL.appspot.com/SECURE_ENDPOINT -H "Authorization: Bearer $TOKEN" \
-H "Destination: slack" \
-H "Space: CBQP83V3M" \
-H "Content-Type: application/json" \
-d '{"text":"yes"}'
```
