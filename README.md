# uptimerobot-operator

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)  [![Go](https://github.com/bennsimon/uptimerobot-operator/actions/workflows/go.yaml/badge.svg?branch=main)](https://github.com/bennsimon/uptimerobot-operator/actions/workflows/go.yaml)

The operator creates, updates, deletes uptimerobot monitors for a particular ingress resource. It's designed to use `friendly_name` attribute of a monitor and/or alert contact for unique identification.

## Description
The operator uses [uptimerobot-tooling](https://github.com/bennsimon/uptimerobot-tooling) to handle api requests.

> The operator will delete the monitor it creates when the ingress resource is deleted.

## Configuration
Environment Variables Supported:

In addition to the environments supplied on the tooling mentioned above, the operator has the following configurations. 

| Variable              | Description                                        | Default     |
|-----------------------|----------------------------------------------------|-------------|
| `DOMAIN_LABEL_PREFIX` | The domain name to use when specifying the labels. | `my.domain` |

With the `DOMAIN_LABEL_PREFIX` as `my.domain` the configurations will be supplied as follows:

````yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: https-minimal-ingress
  labels:
    my.domain/uptimerobot-monitor: "true"
    my.domain/uptimerobot-monitor-type: "HTTP"
    my.domain/uptimerobot-monitor-friendly_name: "tester"
spec:
  rules:
    - host: test-domain.localhost
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: my-service
                port:
                  number: 80
````
The operator reads configurations on the monitor from the label of the ingress resource.

The first label entry `my.domain/uptimerobot-monitor` enables the ingress resource to be evaluated by the operator. The other labels supply the attributes of the monitor. The naming convention is:
`my.domain/uptimerobot-monitor-<attrib>`.

To get more attributes refer to the tooling documentation and uptime robot api documentation.

## Getting Started

You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for
testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever
cluster `kubectl cluster-info` shows).

### Running on the cluster

1. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/uptimerobot-operator:tag
```

2. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/uptimerobot-operator:tag
```

### Uninstall CRDs

To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller

UnDeploy the controller to the cluster:

```sh
make undeploy
```

## Contributing

// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works

This project aims to follow the
Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)
which provides a reconcile function responsible for synchronizing resources until the desired state is reached on the
cluster

### Test It Out

1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions

If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

