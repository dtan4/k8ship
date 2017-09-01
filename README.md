# k8ship

[![Build Status](https://travis-ci.org/dtan4/k8ship.svg?branch=master)](https://travis-ci.org/dtan4/k8ship)
[![codecov](https://codecov.io/gh/dtan4/k8ship/branch/master/graph/badge.svg)](https://codecov.io/gh/dtan4/k8ship)

Deploy image to Kubernetes

## Requirements

- Kubernetes 1.3 or above
- Docker images are tagged with full-qualified Git commit SHA-1 value

## Usage

### `k8ship deploy`

Deploy with Git commit reference (branch name / commit SHA-1 value).
Target deployment and container are detected automatically from Deployment manifest.

To deploy Docker image `dtan4/foo:v3`:

```sh-session
$ k8ship image dtan4/foo:v3
```

### `k8ship image`

Deploy with Docker image.

To deploy Docker image `dtan4/foo:v3` to Deployment `web`:

```sh-session
$ k8ship image dtan4/foo:v3 -d web
```

### `k8ship ref`

Deploy with Git commit reference (branch name | tag | commit SHA-1 value).

To deploy branch `topic/foo` (the latest commit: `fae7c9313f39c382c5051f182bbd281d36368618`) to Deployment `web`:

```sh-session
$ k8ship ref topic/foo -d web
```

Docker image with the tag `fae7c9313f39c382c5051f182bbd281d36368618` will be deployed.

To deploy commit `fae7c93` to Deployment `app`:

```sh-session
$ k8ship ref topic/foo -d app
```

This command does the same as the above.

### `k8ship tag`

Deploy with Docker image tag.

To deploy Docker image tagged `v3` to Deployment `web`:

```sh-session
$ k8ship tag dtan4/foo:v3 -d web
```

## Environment variables

|Key|Description|Required|Example|
|---|---|---|---|
|`GITHUB_ACCESS_TOKEN`|GitHub access token|Required||
|`K8SHIP_ANNOTATION_PREFIX`|Prefix of k8ship-specific annotation|Required|`example.com`|
|`KUBECONFIG`|Path of kubeconfig|||

## Author

Daisuke Fujita ([@dtan4](https://github.com/dtan4))

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
