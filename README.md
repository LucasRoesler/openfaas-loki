# OpenFaaS Loki Provider

## Development flow
OpenFaaS Loki is built with go 1.12+ and uses go modules

### Run unit tets

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


### Build Docker
To enable efficient builds, the Dockerfile uses the experimental RUN syntax to support build-time cahces. This requires enabling the "experimental" features in your Docker installation.  This will enable using [buildkit as the build engine](https://github.com/moby/buildkit/blob/master/frontend/dockerfile/docs/experimental.md#run---mounttypecache)

Build the Docker image and verify the build using

```sh
make build
docker run theaxer/openfaas-loki:latest --version
```

