[![CircleCI](https://circleci.com/gh/devatherock/vela-template-tester.svg?style=svg)](https://circleci.com/gh/devatherock/vela-template-tester)
[![Version](https://img.shields.io/docker/v/devatherock/vela-template-tester?sort=semver)](https://hub.docker.com/r/devatherock/vela-template-tester/)
[![Coverage Status](https://coveralls.io/repos/github/devatherock/vela-template-tester/badge.svg?branch=master)](https://coveralls.io/github/devatherock/vela-template-tester?branch=master)
[![Quality Gate](https://sonarcloud.io/api/project_badges/measure?project=vela-template-tester&metric=alert_status)](https://sonarcloud.io/component_measures?id=vela-template-tester&metric=alert_status&view=list)
[![Docker Pulls](https://img.shields.io/docker/pulls/devatherock/vela-template-tester.svg)](https://hub.docker.com/r/devatherock/vela-template-tester/)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=vela-template-tester&metric=ncloc)](https://sonarcloud.io/component_measures?id=vela-template-tester&metric=ncloc)
[![Docker Image Size](https://img.shields.io/docker/image-size/devatherock/vela-template-tester.svg?sort=date)](https://hub.docker.com/r/devatherock/vela-template-tester/)
# vela-template-tester
API and vela plugin to test and validate [vela-ci templates](https://go-vela.github.io/docs/templates/overview/). The API can also be used to expand any golang/[sprig](https://github.com/Masterminds/sprig) template

## API Reference
### Key parameters:
- **Endpoint**: `https://vela-template-tester.onrender.com/api/expandTemplate`
- **Request Content-Type**: `application/x-yaml`
- **Response Content-Type**: `application/x-yaml`

### Usage samples
**Sample valid template payload:**

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

**Response:**

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

**Sample invalid template payload:**

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

**Response:**

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

**Sample Starlark template payload:**

```
template: |-
  def main(ctx):
    return {
        'version': '1',
        'steps': [
            {
                'name': 'build',
                'image': ctx["vars"]["image"],
                'commands': [
                    'go build',
                    'go test',
                ]
            },
        ],
    }
type: starlark    
parameters:
  image: "go:1.16"
```

## Plugin Reference
### Config
The following parameters can be set to configure the plugin.

**Parameters**
* **input_file** - Input template file to test. Optional if `templates` is specified
* **template_type** - The template type. Needs to be `starlark` if `input_file` is a starlark template
* **variables** - `vars` to test the template with. Doesn't need to be specified if the template can be tested without variables
* **expected_output** - File containing the expected output of the template after applying the variables. Optional, if not specified, only the validity of the processed template will be checked
* **templates** - A list of templates to test. Optional if `input_file` is specified
* **log_level** - Sets the log level. Set to `debug` to enable debug logs. Optional, defaults to `info`

### Examples
**Test a single template**

```yaml
steps:
  - name: vela-template-tester
    ruleset:
      branch: master
      event: [ pull_request, push ]
    image: devatherock/vela-template-tester:latest
    parameters:
      input_file: path/to/template.yml
      variables:
        notification_branch: develop
        notification_event: push
```

**Test a single template with output verification**

```yaml
steps:
  - name: vela-template-tester
    ruleset:
      branch: master
      event: [ pull_request, push ]
    image: devatherock/vela-template-tester:latest
    parameters:
      input_file: path/to/template.yml
      variables:
        notification_branch: develop
        notification_event: push
      expected_output: samples/output_template.yml
```

**Test multiple templates**

```yaml
steps:
  - name: vela-template-tester
    ruleset:
      branch: master
      event: [ pull_request, push ]
    image: devatherock/vela-template-tester:latest
    parameters:
      templates:
        - input_file: path/to/first_template.yml
          variables:
            notification_branch: develop
            notification_event: push
          expected_output: samples/first_template.yml
        - input_file: path/to/second_template.yml
```

## Starlark playground

A vela Starlark template can also be tested using [Starlark playground](https://starpg.onrender.com). We need to specify the template along with the template variables specified within a `ctx` variable and a `print` method call to view the compiled template. Sample usage below:

```
def main(ctx):
  return {
    'version': '1',
    'steps': [
      {
        'name': 'build',
        'image': ctx["vars"]["image"],
        'commands': [
          'go build',
          'go test',
        ]
      },
    ],
}

ctx = {
	"vars": {
		"image": "alpine"
	}
}
print(main(ctx))
```
