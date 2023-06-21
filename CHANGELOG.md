# Changelog

## [Unreleased]
### Changed
- Updated dockerhub readme in CI pipeline
- Switched to `cimg/*` image from `circleci/*` image in CI pipeline
- Set goveralls version to `v0.0.9`, to fix build failure
- Made only HIGH bolt vulnerabilities create issues

## [0.2.0] - 2022-11-05
### Added
- [#2](https://github.com/devatherock/slack-webhook-facade/issues/2): Integrated sonar and coveralls
- Config required to deploy to `render.com`

### Removed
- Deployment to heroku

## [0.1.0] - 2020-06-19
### Added
- Initial version. Exposes an endpoint that accepts a slack webhook payload and posts it to Zulip