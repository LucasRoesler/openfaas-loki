# OpenFaaS Loki

A [Loki](https://github.com/grafana/loki) powered log provider for [OpenFaaS](https://www.openfaas.com)

[![Go Report Card](https://goreportcard.com/badge/github.com/LucasRoesler/openfaas-loki)](https://goreportcard.com/report/github.com/LucasRoesler/openfaas-loki) ![GitHub branch checks state](https://img.shields.io/github/checks-status/LucasRoesler/openfaas-loki/main)

OpenFaaS Loki implementes the new log provider interface from [`faas-provider`](https://github.com/openfaas/faas-provider). This means you can query Loki for your function logs with the [`faas-cli`](https://github.com/openfaas/faas-cli)!

_Limitations_: This initial version does not support log tail streams. This means that `faas-cli logs` will alway behave as if `--follow=false`. Live tailing of the logs is currently being worked on.

All other OpenFaaS log behaviors should be fully supported. [Issues are welcome](https://github.com/LucasRoesler/openfaas-loki/issues/new) if you notice a bug or room for improvement.

## Install with Helm

1. Clone this repo
2. Determine the URL of your Loki installation, if you used the default Helm values to install Loki, this will be `http://<service name>.<namespace>:3100`
3. Then install the `openfaas-loki` provider using Helm:
   ```sh
   helm upgrade --install ofloki oci:://ghcr.io/lucasroesler/charts/openfaas-loki \
       --namespace openfaas \
       --set lokiURL=http://loki.monitoring:3100 \
       --set logLevel=DEBUG
   ```
4. Then update the `gateway` with the environment variable described in the NOTES output of the helm install. Currently, this can be done using
   ```sh
   kubectl -n openfaas set env deployment/gateway -c gateway -e logs_provider_url=http://ofloki-openfaas-loki.openfaas:9191/
   ```
   The environment variable is only needed on the `gateway` container.
5. Test the installation using
   ```sh
   faas-cli store deploy nodeinfo
   echo "" | faas-cli invoke nodeinfo
   faas-cli logs nodeinfo --tail=3
   ```

## Development flow

OpenFaaS Loki is built with go 1.12+ and uses go modules and all of the development actions are currently defined through the `Makefile`

### Run tests

```sh
make test
```

### Local install

You can install using

```sh
make install
```

Check the installation using

```sh
openfaas-loki --version
```

### Docker Build

Build the Docker image for amd64 archs and verify the build using

```sh
make image ARCH=linux/amd64 OWNER=lucasroesler REGISTRY=ghcr.io
docker run ghcr.io/lucasroesler/openfaas-loki:dev --version
```
