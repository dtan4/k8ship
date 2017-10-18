# k8ship

[![Build Status](https://travis-ci.org/dtan4/k8ship.svg?branch=master)](https://travis-ci.org/dtan4/k8ship)
[![codecov](https://codecov.io/gh/dtan4/k8ship/branch/master/graph/badge.svg)](https://codecov.io/gh/dtan4/k8ship)

Deploy image to Kubernetes

## Requirements

Kubernetes 1.3 or above

## Application Preparation

### Docker image

Docker images are tagged with full-qualified Git commit SHA-1 value

e.g. `quay.io/dtan4/k8ship:0118ef0b66a6b9cb04a6547aca5a17d0ad601782`

### Kubernetes Deployment

To use k8ship deploy, you have to add a few annotations to your Deployment manifest.

|Key|Description|
|---|---|
|`example.com/deploy-target`|`"true"/"false"` whether this Deployment can be deployed by k8ship|
|`example.com/deploy-target-container`|Container name which will be updated by k8ship|
|`example.com/github`|Pair of the target container and its GitHub repository. `<container>=<user>/<repo>`|

NOTE: The prefix `example.com` can be replaced as you like via `K8SHIP_ANNOTATION_PREFIX`.

You MUST add `example.com/deploy-target="true"` to deploy the Deployment using k8ship.

#### 1 Pod, 1 Container

If there is 1 container in the target Pod, all you need to do is that specifying GitHub repository by `example.com/github` annotation.

Following manifest shows that `web` container will be deployed from `dtan4/awesome-app` repository.

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: awesome-app
  labels:
    name: awesome-app
    role: web
  annotations:
    dtan4.net/deploy-target: "true"          # <===== ADDED
    dtan4.net/github: web=dtan4/awesome-app  # <===== ADDED
spec:
  replicas: 1
  template:
    spec:
      containers:
      - image: quay.io/dtan4/awesome-app:latest
        name: web
```

#### 1 Pod, N Containers

Let us assume that there is a Pod including web application and Nginx.

k8ship can deploy one container per one Pod.
You have to specify which container is k8ship deploy target.

Following manifest shows that `web` container will be deployed from `dtan4/awesome-app` repository.
`nginx` container will not be updated by k8ship.

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: awesome-app
  labels:
    name: awesome-app
    role: web
  annotations:
    dtan4.net/deploy-target: "true"          # <===== ADDED
    dtan4.net/deploy-target-container: web   # <===== ADDED
    dtan4.net/github: web=dtan4/awesome-app  # <===== ADDED
spec:
  replicas: 1
  template:
    spec:
      containers:
      - image: quay.io/dtan4/awesome-app:latest
        name: web
      - image: nginx:latest
        name: nginx
```

## Command-line Usage

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
