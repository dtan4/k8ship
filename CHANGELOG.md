# [v0.7.1](https://github.com/dtan4/k8ship/releases/tag/v0.7.1) (2018-01-25)

## Features

- Change Deployment status: pending -> success [#30](https://github.com/dtan4/k8ship/pull/30)

# [v0.7.0](https://github.com/dtan4/k8ship/releases/tag/v0.7.0) (2018-01-24)

## Features

- Create GitHub Deployment at `k8ship deploy` [#28](https://github.com/dtan4/k8ship/pull/28)

# [v0.6.0](https://github.com/dtan4/k8ship/releases/tag/v0.6.0) (2017-12-18)

## Features

- Replace "DEPLOYED IMAGE" to "IMAGE" in header [#27](https://github.com/dtan4/k8ship/pull/27)
- Attach user name with reload [#25](https://github.com/dtan4/k8ship/pull/25)
- Print deploy user in `k8ship history` output [#24](https://github.com/dtan4/k8ship/pull/24)
- Check history length before extracting [#23](https://github.com/dtan4/k8ship/pull/2)3
- Deploy with username [#22](https://github.com/dtan4/k8ship/pull/22)
- Add `k8ship history` command [#21](https://github.com/dtan4/k8ship/pull/21)

# [v0.5.0](https://github.com/dtan4/k8ship/releases/tag/v0.5.0) (2017-12-05)

## Features

- `k8ship reload` reloads all Pods in Deployment [#18](https://github.com/dtan4/k8ship/pull/18)

# [v0.4.0](https://github.com/dtan4/k8ship/releases/tag/v0.4.0) (2017-11-27)

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
