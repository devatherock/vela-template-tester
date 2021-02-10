[![CircleCI](https://circleci.com/gh/devatherock/vela-template-tester.svg?style=svg)](https://circleci.com/gh/devatherock/vela-template-tester)
[![Version](https://img.shields.io/docker/v/devatherock/vela-template-tester?sort=semver)](https://hub.docker.com/r/devatherock/vela-template-tester/)
[![Coverage Status](https://coveralls.io/repos/github/devatherock/vela-template-tester/badge.svg?branch=master)](https://coveralls.io/github/devatherock/vela-template-tester?branch=master)
[![Quality Gate](https://sonarcloud.io/api/project_badges/measure?project=vela-template-tester&metric=alert_status)](https://sonarcloud.io/component_measures?id=vela-template-tester&metric=alert_status&view=list)
[![Docker Pulls](https://img.shields.io/docker/pulls/devatherock/vela-template-tester.svg)](https://hub.docker.com/r/devatherock/vela-template-tester/)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=vela-template-tester&metric=ncloc)](https://sonarcloud.io/component_measures?id=vela-template-tester&metric=ncloc)
[![Docker Image Size](https://img.shields.io/docker/image-size/devatherock/vela-template-tester.svg?sort=date)](https://hub.docker.com/r/devatherock/vela-template-tester/)
# vela-template-tester
API and vela plugin to test and validate [vela-ci templates](https://go-vela.github.io/docs/templates/overview/)

## API Reference
### Key parameters:
- **Endpoint**: `https://vela-template-tester.herokuapp.com/api/expandTemplate`
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

## Plugin Reference
Please refer [docs](DOCS.md)
