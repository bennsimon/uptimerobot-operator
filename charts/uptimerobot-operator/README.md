
Uptimerobot-operator
===========

This operator creates, updates, deletes uptimerobot monitors for a particular ingress resource. It's designed to use friendly_name attribute of a monitor and/or alert contact for unique identification.


## TL;DR

```bash
$ helm repo add uptimerobot-operator https://bennsimon.github.io/uptimerobot-operator/
$ helm install uptimerobot-operator uptimerobot-operator/uptimerobot-operator
```

## Introduction

This chart bootstraps  [uptimerobot-operator](https://github.com/bennsimon/uptimerobot-operator) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.12+
- Helm 3.1.0

## Installing the Chart

To install the chart with the release name `uptimerobot-operator`:

## Configuration

The following table lists the configurable parameters of the Uptimerobot-operator chart and their default values.

| Parameter                                    | Description | Default                                                    |
|----------------------------------------------|-------------|------------------------------------------------------------|
| `replicaCount`                               |             | `1`                                                        |
| `image.repository`                           |             | `"bennsimon/uptimerobot-operator"`                         |
| `image.pullPolicy`                           |             | `"IfNotPresent"`                                           |
| `image.tag`                                  |             | `"v0.0.2"`                                        |
| `imagePullSecrets`                           |             | `[]`                                                       |
| `nameOverride`                               |             | `""`                                                       |
| `fullnameOverride`                           |             | `""`                                                       |
| `clusterRole.create`                         |             | `true`                                                     |
| `serviceAccount.create`                      |             | `true`                                                     |
| `serviceAccount.annotations`                 |             | `{}`                                                       |
| `serviceAccount.name`                        |             | `""`                                                       |
| `podAnnotations`                             |             | `{}`                                                       |
| `podSecurityContext`                         |             | `{}`                                                       |
| `securityContext.allowPrivilegeEscalation`   |             | `false`                                                    |
| `securityContext.capabilities.drop`          |             | `["ALL"]`                                                  |
| `resources.limits.memory`                    |             | `"128Mi"`                                                  |
| `resources.requests.cpu`                     |             | `"10m"`                                                    |
| `resources.requests.memory`                  |             | `"64Mi"`                                                   |
| `autoscaling.enabled`                        |             | `false`                                                    |
| `autoscaling.minReplicas`                    |             | `1`                                                        |
| `autoscaling.maxReplicas`                    |             | `100`                                                      |
| `autoscaling.targetCPUUtilizationPercentage` |             | `80`                                                       |
| `nodeSelector`                               |             | `{}`                                                       |
| `tolerations`                                |             | `[]`                                                       |
| `affinity`                                   |             | `{}`                                                       |
| `selectorLabels.control-plane`               |             | `"controller-manager"`                                     |
| `livenessProbe.httpGet.path`                 |             | `"/healthz"`                                               |
| `livenessProbe.httpGet.port`                 |             | `8081`                                                     |
| `livenessProbe.initialDelaySeconds`          |             | `15`                                                       |
| `livenessProbe.periodSeconds`                |             | `20`                                                       |
| `readinessProbe.httpGet.path`                |             | `"/readyz"`                                                |
| `readinessProbe.httpGet.port`                |             | `8081`                                                     |
| `readinessProbe.initialDelaySeconds`         |             | `5`                                                        |
| `readinessProbe.periodSeconds`               |             | `10`                                                       |
| `env`                                        |             | `[{"name": "UPTIME_ROBOT_API_KEY", "value": "<api-key>"}]` |
