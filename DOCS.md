## Config
The following parameters can be set to configure the plugin.

### Parameters
* **input_file** - Input template file to test. Optional if `templates` is specified
* **variables** - `vars` to test the template with. Doesn't need to be specified if the template can be tested without variables
* **expected_output** - File containing the expected output of the template after applying the variables. Optional, if not specified, only the validity of the processed template will be checked
* **templates** - A list of templates to test. Optional if `input_file` is specified
* **log_level** - Sets the log level. Set to `debug` to enable debug logs. Optional, defaults to `info`

## Examples
### Test a single template

```yaml
steps:
  - name: vela-template-tester
    ruleset:
      branch: master
      event: [ pull_request, push ]
    image: devatherock/vela-template-tester:0.2.0
    parameters:
      input_file: path/to/template.yml
      variables:
        notification_branch: develop
        notification_event: push
```

### Test a single template with output verification

```yaml
steps:
  - name: vela-template-tester
    ruleset:
      branch: master
      event: [ pull_request, push ]
    image: devatherock/vela-template-tester:0.2.0
    parameters:
      input_file: path/to/template.yml
      variables:
        notification_branch: develop
        notification_event: push
      expected_output: samples/output_template.yml
```

### Test multiple templates

```yaml
steps:
  - name: vela-template-tester
    ruleset:
      branch: master
      event: [ pull_request, push ]
    image: devatherock/vela-template-tester:0.2.0
    parameters:
      templates:
        - input_file: path/to/first_template.yml
          variables:
            notification_branch: develop
            notification_event: push
          expected_output: samples/first_template.yml
        - input_file: path/to/second_template.yml
```
