# Changelog

## [Unreleased]
### Added
- `circleci-templates` orb for common tasks

### Changed
- Made only HIGH bolt vulnerabilities create issues
- fix(deps): update module github.com/stretchr/testify to v1.9.0
- fix(deps): update module github.com/urfave/cli/v2 to v2.27.4
- chore(deps): bump golang.org/x/crypto from 0.14.0 to 0.17.0
- Upgraded `golang.org/x/net` to `0.23.0`
- fix(deps): update go.starlark.net digest to 42030a7
- chore(deps): update alpine docker tag to v3.20.3
- fix(deps): update module github.com/masterminds/sprig/v3 to v3.3.0
- Upgraded go to `1.23`
- fix(deps): update go.starlark.net digest to 1207426
- fix(deps): update module github.com/urfave/cli/v2 to v2.27.5
- fix(deps): update go.starlark.net digest to d4d7611
- fix(deps): update module github.com/stretchr/testify to v1.10.0
- fix(deps): update go.starlark.net digest to c05ff20
- chore(deps): update alpine docker tag to v3.21.0
- fix(deps): update go.starlark.net digest to 8dfa5b9
- chore(deps): bump golang.org/x/crypto from 0.26.0 to 0.31.0
- chore(deps): update alpine docker tag to v3.21.2
- fix(deps): update go.starlark.net digest to d908c3e
- fix(deps): update go.starlark.net digest to 492d367
- chore(deps): update alpine docker tag to v3.21.3
- fix(deps): update go.starlark.net digest to 0d3f41d
- fix(deps): update module github.com/urfave/cli/v2 to v2.27.6
- fix(deps): update go.starlark.net digest to d9371fe

### Removed
- Trailing spaces in yaml files
- Dependence on `starpg`

## [1.0.0] - 2023-06-16
### Added
- [#13](https://github.com/devatherock/vela-template-tester/issues/13): Functional tests
- `changelog-updater` to CI pipeline

### Changed
- Combined jobs in CI pipeline using parameters
- [#42](https://github.com/devatherock/vela-template-tester/issues/42): Merged contents of `DOCS.md` into `README.md`
- [#40](https://github.com/devatherock/vela-template-tester/issues/40): Updated dockerhub readme in CI pipeline
- Restructured the project as per [golang-standards](https://github.com/golang-standards/project-layout)
- [#31](https://github.com/devatherock/vela-template-tester/issues/31): Upgraded `goutils` to `1.1.1`
- [#29](https://github.com/devatherock/vela-template-tester/issues/29): Upgraded `testify` to `1.8.4`
- Upgraded go to `1.20`
- Upgraded `sprig` to `3.2.3`
- [#47](https://github.com/devatherock/vela-template-tester/issues/47): Upgraded `logrus` to `1.9.3`
- [#41](https://github.com/devatherock/vela-template-tester/issues/41): Built a multi-arch docker image
- Upgraded `cli` to `2.25.6`
- Upgraded `alpine` to `3.18.2`

### Removed
- Unused `PORT` environment variable from render

## [0.5.0] - 2022-08-27
### Added
- Deployment to render.com

### Changed
- Used `starpg` deployed on render.com

### Removed
- Deployment to heroku

## [0.4.0] - 2021-08-21
### Added
- Support for Starlark based vela templates

## [0.3.1] - 2021-06-23
### Changed
- [#25](https://github.com/devatherock/vela-template-tester/issues/25): Fixed failures in templates that used `vela` function

## [0.3.0] - 2021-04-15
### Added
- test: Basic tests so that coveralls and Sonar can be introduced
- Accepted a list as `parameters` in the API along with `map`. This will make the API be able to expand any golang/sprig template, not just vela templates

## [0.2.0] - 2020-09-26
### Added
- vela plugin to test vela templates

## [0.1.3] - 2020-06-06
### Added
- Ability to use `PORT` environment variable as port, for Heroku.

## [0.1.2] - 2020-06-06
### Added
- Step to deploy the API to heroku

## [0.1.1] - 2020-05-30
### Added
- A health check

## [0.1.0] - 2020-05-30
### Added
- API to test and validate a vela-ci template