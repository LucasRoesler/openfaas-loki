# Changelog

## [1.7.3](https://github.com/LucasRoesler/openfaas-loki/compare/v1.7.2...v1.7.3) (2023-09-23)


### Bug Fixes

* use the releaser token instead of the default ([0cefe7c](https://github.com/LucasRoesler/openfaas-loki/commit/0cefe7c761c77eabb46ab0698e876b53ac7f125a))

## [1.7.2](https://github.com/LucasRoesler/openfaas-loki/compare/v1.7.1...v1.7.2) (2023-09-23)


### Miscellaneous

* release 1.7.2 ([68798ff](https://github.com/LucasRoesler/openfaas-loki/commit/68798ffdd3db3388256247d989b384bbf39f05c7))

## [1.7.1](https://github.com/LucasRoesler/openfaas-loki/compare/v1.7.0...v1.7.1) (2023-09-23)


### Bug Fixes

* handle go 1.18 amd versioning during builds ([3cc3b33](https://github.com/LucasRoesler/openfaas-loki/commit/3cc3b3318e9567f455b4fb47e4ee39c5d1405252))

## [1.7.0](https://github.com/LucasRoesler/openfaas-loki/compare/v1.6.0...v1.7.0) (2023-09-23)


### Features

* upgrade to go 1.21 and switch to slog ([a377237](https://github.com/LucasRoesler/openfaas-loki/commit/a377237dfb0dd8e1a29319de7b614248dd1cbbf9))


### Bug Fixes

* update dev tools and fix linting errors ([7ce6849](https://github.com/LucasRoesler/openfaas-loki/commit/7ce684961f4b54b67ab91506e8321ab2314d1292))


### Miscellaneous

* bump azure/setup-helm from 1 to 3 ([#40](https://github.com/LucasRoesler/openfaas-loki/issues/40)) ([0a76b00](https://github.com/LucasRoesler/openfaas-loki/commit/0a76b008e30c576968043899010f58f8318de212))
* bump docker/build-push-action from 2 to 5 ([#61](https://github.com/LucasRoesler/openfaas-loki/issues/61)) ([4eff495](https://github.com/LucasRoesler/openfaas-loki/commit/4eff4951d2d1b155ee1f1b2d2964e17cf6362690))
* bump docker/login-action from 1 to 3 ([#59](https://github.com/LucasRoesler/openfaas-loki/issues/59)) ([6f8590a](https://github.com/LucasRoesler/openfaas-loki/commit/6f8590a8c72e065c29474449400e81be5ffd67d8))
* bump docker/metadata-action from 3 to 5 ([#60](https://github.com/LucasRoesler/openfaas-loki/issues/60)) ([0720400](https://github.com/LucasRoesler/openfaas-loki/commit/072040005c41d46e0f9e721171f9c52db18b5909))
* bump docker/setup-qemu-action from 1 to 3 ([#58](https://github.com/LucasRoesler/openfaas-loki/issues/58)) ([46acb42](https://github.com/LucasRoesler/openfaas-loki/commit/46acb425feb20b024843f46f79980bac76374bec))
* Update repo add instructions ([e5aeb8c](https://github.com/LucasRoesler/openfaas-loki/commit/e5aeb8c8b2c5afd24665a430debdfcb1c195ca21))


### Automations

* add missing permission for release-please ([68ea751](https://github.com/LucasRoesler/openfaas-loki/commit/68ea7516feb037a8169634e353cb713e133f1826))
* bump actions/checkout from 2 to 4 ([#64](https://github.com/LucasRoesler/openfaas-loki/issues/64)) ([a15f6c9](https://github.com/LucasRoesler/openfaas-loki/commit/a15f6c93e8d62a173be3f1990bb5fab1efdcf70d))
* bump actions/setup-go from 3 to 4 ([#63](https://github.com/LucasRoesler/openfaas-loki/issues/63)) ([5146af5](https://github.com/LucasRoesler/openfaas-loki/commit/5146af58540beb48a8e2b4022db0485db2904612))
* bump docker/setup-buildx-action from 1 to 3 ([#65](https://github.com/LucasRoesler/openfaas-loki/issues/65)) ([14c6931](https://github.com/LucasRoesler/openfaas-loki/commit/14c6931d8a43eeb82d8277c8b087bba4a137edc9))
* bump goreleaser/goreleaser-action from 2 to 5 ([#66](https://github.com/LucasRoesler/openfaas-loki/issues/66)) ([f250a50](https://github.com/LucasRoesler/openfaas-loki/commit/f250a5032e36685f2e8ac7f1461a2f66806fc194))
* fix golangci-lint-action version ([8086e23](https://github.com/LucasRoesler/openfaas-loki/commit/8086e232ff5435e8f23d88b3d0babd9db327f3d5))
* update dependabot settings for actions ([3e879b0](https://github.com/LucasRoesler/openfaas-loki/commit/3e879b0b557c316d5cb67a2cc183f17a4eee5270))
* update release flow to push OCI helm charts ([7492502](https://github.com/LucasRoesler/openfaas-loki/commit/7492502ae8a0de7b7004ac8625a43fc7b9177ad2))

## [1.6.0](https://github.com/LucasRoesler/openfaas-loki/compare/v1.5.0...v1.6.0) (2022-03-08)


### Features

* Add annotations to the helm chart ([#6](https://github.com/LucasRoesler/openfaas-loki/issues/6)) ([9c1cf32](https://github.com/LucasRoesler/openfaas-loki/commit/9c1cf32019aab44f20e57862cbc74285f6723796))
* add request logging middleware ([#20](https://github.com/LucasRoesler/openfaas-loki/issues/20)) ([ea0ea35](https://github.com/LucasRoesler/openfaas-loki/commit/ea0ea350c1d0358a404cbc4e8fa76a724800aba8))


### Bug Fixes

* include v prefix for the appVersion ([39ee398](https://github.com/LucasRoesler/openfaas-loki/commit/39ee398777eb116d6cc1ca2a650673251d868c0c))


### Miscellaneous

* autoupdate dependencies ([#17](https://github.com/LucasRoesler/openfaas-loki/issues/17)) ([3f03bea](https://github.com/LucasRoesler/openfaas-loki/commit/3f03beac2fba85ee63c6ffbfc66144bf4e36800f)), closes [#15](https://github.com/LucasRoesler/openfaas-loki/issues/15) [#14](https://github.com/LucasRoesler/openfaas-loki/issues/14) [#13](https://github.com/LucasRoesler/openfaas-loki/issues/13) [#12](https://github.com/LucasRoesler/openfaas-loki/issues/12) [#11](https://github.com/LucasRoesler/openfaas-loki/issues/11)
* bump actions/checkout from 2 to 3 ([#9](https://github.com/LucasRoesler/openfaas-loki/issues/9)) ([b79d098](https://github.com/LucasRoesler/openfaas-loki/commit/b79d0982ca936b7f467876588f01138bca1be1be))
* bump release-please to v3 ([f651ffe](https://github.com/LucasRoesler/openfaas-loki/commit/f651ffe6b70f08daa50be12277c6fd0ea41a569f))
* rename chart folder to prepare for self-hosting ([1f8325c](https://github.com/LucasRoesler/openfaas-loki/commit/1f8325c28dce34c9cf7930b443dc529c36492be6))
* update build flows and add github actions CI ([#7](https://github.com/LucasRoesler/openfaas-loki/issues/7)) ([04b2b44](https://github.com/LucasRoesler/openfaas-loki/commit/04b2b4424f880745fb0c264016668ad8c81a1dc4))
* upgrade loki to the commit for v2.4.2 ([#18](https://github.com/LucasRoesler/openfaas-loki/issues/18)) ([c37c4fd](https://github.com/LucasRoesler/openfaas-loki/commit/c37c4fd70cd2114741569efac1c0fdb99fff1573))
* use appVersion as the default image tag in Helm ([60d09b2](https://github.com/LucasRoesler/openfaas-loki/commit/60d09b2ac502ba57e1b66fffc68c22082c3115e8))


### Automations

* add ci changes to the release notes ([830c047](https://github.com/LucasRoesler/openfaas-loki/commit/830c047a1fec23fd0cfe16505bde40aa0beddff4))
* add helm chart lint and test flow ([b3ab439](https://github.com/LucasRoesler/openfaas-loki/commit/b3ab439a9e6ce2f07a01d1253adeaf1a093b84c5))
* improve release automation ([#19](https://github.com/LucasRoesler/openfaas-loki/issues/19)) ([7b79029](https://github.com/LucasRoesler/openfaas-loki/commit/7b7902984e0b306c901dc5006540dc49e8dd4f27))
* Upload helm packages during release ([43613e4](https://github.com/LucasRoesler/openfaas-loki/commit/43613e4d85a8c60f692bcaea31ff6781f37c9018))
