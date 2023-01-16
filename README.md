# uptimerobot-operator

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)  [![Go](https://github.com/bennsimon/uptimerobot-operator/actions/workflows/go.yaml/badge.svg?branch=main)](https://github.com/bennsimon/uptimerobot-operator/actions/workflows/go.yaml)

The operator creates, updates, deletes uptimerobot monitors for a particular ingress resource. It's designed to use `friendly_name` attribute of a monitor and/or alert contact for unique identification.

## Description

The operator uses [uptimerobot-tooling](https://github.com/bennsimon/uptimerobot-tooling) to handle api requests.

> The operator will delete the monitor it creates when the ingress resource is deleted.

## Configuration

Environment Variables Supported:

In addition to the environments supplied on the tooling mentioned above, the operator has the following configurations.

| Variable        | Description                                             | Default               |
|-----------------|---------------------------------------------------------|-----------------------|
| `DOMAIN_PREFIX` | The domain name to use when specifying the annotations. | `bennsimon.github.io` |

With the `DOMAIN_PREFIX` as `bennsimon.github.io` the configurations will be supplied as follows:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: https-minimal-ingress
  annotations:
      bennsimon.github.io/uptimerobot-monitor: "true"
      bennsimon.github.io/uptimerobot-monitor-type: "HTTP"
      bennsimon.github.io/uptimerobot-monitor-friendly_name: "tester"
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
```

The operator reads configurations of the monitor from the annotation on the ingress resource.

The first annotation entry `bennsimon.github.io/uptimerobot-monitor` enables the ingress resource to be evaluated by the operator. The other annotations supply the attributes of the monitor. The naming convention is:
`bennsimon.github.io/uptimerobot-monitor-<attrib>`.

To get more attributes refer to the tooling documentation and uptime robot api documentation.

## Getting Started

You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for
testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever
cluster `kubectl cluster-info` shows).

> Ensure you have supplied the environment variables here config/manager/manager.yaml.

### Deploying on the cluster

To deploy the operator you will need the following manifests:

*   serviceaccount
*   clusterrole
*   clusterrolebinding
*   deployment

Create a yaml file and paste the yaml snippet below and update configurations to your preferences.

```yaml (uptimerobot-operator.yaml)
---
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app.kubernetes.io/name: serviceaccount
    app.kubernetes.io/instance: controller-manager
    app.kubernetes.io/component: rbac
    app.kubernetes.io/created-by: uptimerobot-operator
    app.kubernetes.io/part-of: uptimerobot-operator
  name: uptimerobot-operator
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: uptimerobot-operator
rules:
- apiGroups:
  - networking.k8s.io
  resources:
  - ingresses
  verbs:
  - get
  - list
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  labels:
    app.kubernetes.io/name: clusterrolebinding
    app.kubernetes.io/instance: manager-rolebinding
    app.kubernetes.io/component: rbac
    app.kubernetes.io/created-by: uptimerobot-operator
    app.kubernetes.io/part-of: uptimerobot-operator
  name: uptimerobot-operator
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: uptimerobot-operator
subjects:
  - kind: ServiceAccount
    name: uptimerobot-operator
    namespace: default # update this to preferred namespace
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: uptimerobot-operator
  labels:
    control-plane: controller-manager
    app.kubernetes.io/name: deployment
    app.kubernetes.io/instance: controller-manager
    app.kubernetes.io/component: manager
    app.kubernetes.io/created-by: uptimerobot-operator
    app.kubernetes.io/part-of: uptimerobot-operator
spec:
  selector:
    matchLabels:
      control-plane: controller-manager
  replicas: 1
  template:
    metadata:
      annotations:
        kubectl.kubernetes.io/default-container: manager
      labels:
        control-plane: controller-manager
    spec:
      containers:
        - command:
            - /manager
          env:
            - name: UPTIME_ROBOT_API_KEY
              value: "<api-key>"
#            - name: MONITOR_RESOLVE_ALERT_CONTACTS_BY_FRIENDLY_NAME
#              value: "true"
#            - name: MONITOR_ALERT_CONTACTS_DELIMITER
#              value: "-"
#            - name: MONITOR_ALERT_CONTACTS_ATTRIB_DELIMITER
#              value: "_"
          image: bennsimon/uptimerobot-operator:v0.1.0
          name: manager
          securityContext:
            allowPrivilegeEscalation: false
            capabilities:
              drop:
                - "ALL"
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8081
            initialDelaySeconds: 15
            periodSeconds: 20
          readinessProbe:
            httpGet:
              path: /readyz
              port: 8081
            initialDelaySeconds: 5
            periodSeconds: 10
          # TODO(user): Configure the resources accordingly based on the project requirements.
          # More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
          resources:
            limits:
              cpu: 500m
              memory: 128Mi
            requests:
              cpu: 10m
              memory: 64Mi
      serviceAccountName: uptimerobot-operator
      terminationGracePeriodSeconds: 10
```

    kubectl apply -f uptimerobot-operator.yaml

## Development

### Running on the cluster

Use the latest tag from [dockerhub](https://hub.docker.com/r/bennsimon/uptimerobot-operator/tags)

1.  Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=bennsimon/uptimerobot-operator:v0.1.0
```

2.  Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=bennsimon/uptimerobot-operator:v0.1.0
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

1.  Install the CRDs into the cluster:

```sh
make install
```

2.  Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

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
