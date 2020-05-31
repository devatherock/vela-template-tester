[![CircleCI](https://circleci.com/gh/devatherock/vela-template-tester.svg?style=svg)](https://circleci.com/gh/devatherock/vela-template-tester)
[![Docker Pulls](https://img.shields.io/docker/pulls/devatherock/vela-template-tester-api.svg)](https://hub.docker.com/r/devatherock/vela-template-tester-api/)
[![Docker Image Size](https://img.shields.io/docker/image-size/devatherock/vela-template-tester-api.svg?sort=date)](https://hub.docker.com/r/devatherock/vela-template-tester-api/)
[![Docker Image Layers](https://img.shields.io/microbadger/layers/devatherock/vela-template-tester-api.svg)](https://microbadger.com/images/devatherock/vela-template-tester-api)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
# vela-template-tester
API to test and validate vela-ci templates

## API Reference
### Key parameters:
- **Endpoint**: `${host}/api/expandTemplate`
- **Request Content-Type**: `application/x-yaml`
- **Response Content-Type**: `application/x-yaml`

### Usage samples
#### Sample valid template payload:

```yaml
template: |-
  metadata:
    template: true

  slack_plugin_image: &slack_plugin_image
    image: devatherock/simple-slack:0.2.0

  steps:
    - name: notify_success
      ruleset:
        branch: {{ default "[ master, v1 ]" .notification_branch }}
        event: {{ default "[ push, tag ]" .notification_event }}
      <<: *slack_plugin_image
      secrets: [ slack_webhook ]
      parameters:
        color: "#33ad7f"
        text: |-
          Success: {{"{{.BuildLink}}"}} ({{"{{.BuildRef}}"}}) by {{"{{.BuildAuthor}}"}}
          {{"{{.BuildMessage}}"}}
parameters:
  notification_branch: develop
```

#### Response:

```yaml
message: template is a valid yaml
template: |-
  metadata:
    template: true

  slack_plugin_image: &slack_plugin_image
    image: devatherock/simple-slack:0.2.0

  steps:
    - name: notify_success
      ruleset:
        branch: develop
        event: [ push, tag ]
      <<: *slack_plugin_image
      secrets: [ slack_webhook ]
      parameters:
        color: "#33ad7f"
        text: |-
          Success: {{.BuildLink}} ({{.BuildRef}}) by {{.BuildAuthor}}
          {{.BuildMessage}}
```

#### Sample invalid template payload:

```yaml
template: |-
  metadata:
    template: true

  steps:
    - name: notify_success
      ruleset:
        branch: {{ default "[ master, v1 ]" .notification_branch }}
        event: {{ default "[ push, tag ]" .notification_event }}
      <<: *slack_plugin_image
      secrets: [ slack_webhook ]
      parameters:
        color: "#33ad7f"
        text: |-
          Success: {{"{{.BuildLink}}"}} ({{"{{.BuildRef}}"}}) by {{"{{.BuildAuthor}}"}}
          {{"{{.BuildMessage}}"}}
parameters:
  notification_branch: develop
```

#### Response:

```yaml
message: template is not a valid yaml
error: 'yaml: unknown anchor ''slack_plugin_image'' referenced'
template: |-
  metadata:
    template: true

  steps:
    - name: notify_success
      ruleset:
        branch: develop
        event: [ push, tag ]
      <<: *slack_plugin_image
      secrets: [ slack_webhook ]
      parameters:
        color: "#33ad7f"
        text: |-
          Success: {{.BuildLink}} ({{.BuildRef}}) by {{.BuildAuthor}}
          {{.BuildMessage}}
```