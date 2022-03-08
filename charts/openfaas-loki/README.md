# openfaas-loki

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![AppVersion: v1.5.0](https://img.shields.io/badge/AppVersion-v1.5.0-informational?style=flat-square)

A Loki powered log provider for OpenFaaS

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Lucas Roesler |  | https://lucasroesler.com |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| annotations | object | `{}` |  |
| image.pullPolicy | string | `"Always"` |  |
| image.repository | string | `"ghcr.io/lucasroesler/openfaas-loki"` |  |
| image.tag | string | `""` | will default to Chart.appVersion |
| logLevel | string | `"INFO"` |  |
| lokiURL | string | `""` | required |
| nodeSelector | object | `{}` |  |
| replicaCount | int | `1` |  |
| resources.requests.cpu | string | `"50m"` |  |
| resources.requests.memory | string | `"120Mi"` |  |
| service.port | int | `9191` |  |
| service.type | string | `"ClusterIP"` |  |
| timeout | string | `"30s"` |  |
| tolerations | list | `[]` |  |

