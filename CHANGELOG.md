# [v0.4.0](https://github.com/dtan4/k8ship/releases/tag/v0.3.0) (2017-11-27)

## Features

- `k8ship deploy` does not fail if the target image _tags_ are different [#16](https://github.com/dtan4/k8ship/pull/16)
  - If the target images have the same image name but do not have the same image tag (e.g.,  `web:latest` and `web:abc123`), `k8ship deploy` does not fail and replaces to the given image.
  - If the target images do not have the same image name (e.g., `web:latest` and `nginx:latest`), `k8ship deploy` will fail.

# [v0.3.0](https://github.com/dtan4/k8ship/releases/tag/v0.3.0) (2017-07-19)

## Breaking changes

- Annotation nae has been changed
  - Remove / from annotation [#14](https://github.com/dtan4/k8ship/pull/14)
  - `<prefix>deploy/target` -> `<prefix>deploy-target`
  - `<prefix>deploy/target-container` -> `<prefix>deploy-target-container`

# [v0.2.1](https://github.com/dtan4/k8ship/releases/tag/v0.2.2) (2017-07-19)

## Features

- Add annotation prefix [#12](https://github.com/dtan4/k8ship/pull/12)

# [v0.2.1](https://github.com/dtan4/k8ship/releases/tag/v0.2.1) (2017-07-19)

## Features

- Print finish message [#10](https://github.com/dtan4/k8ship/pull/10)
- Add `k8ship deploy --tag` and `k8ship deploy --image` [#9](https://github.com/dtan4/k8ship/pull/9)

# [v0.2.0](https://github.com/dtan4/k8ship/releases/tag/v0.2.0) (2017-07-18)

## Features

- Set change cause for history [#7](https://github.com/dtan4/k8ship/pull/7)
- Update image using partial patch [#6](https://github.com/dtan4/k8ship/pull/6)
- Deploy to multiple Deployments at once (new `kube deploy`) [#5](https://github.com/dtan4/k8ship/pull/5)
- Separate (old) `k8ship deploy` to `k8ship tag` and `k8ship image` [#3](https://github.com/dtan4/k8ship/pull/3)

# [v0.1.0](https://github.com/dtan4/k8ship/releases/tag/v0.1.0) (2017-07-14)

Initial release.
