# .section-proxy.yaml

The `.section-proxy.yaml` file is a YAML-formatted file with a custom schema for
describing how a proxy will operate on the section.io platform.

This file should be located in the same directory as a proxy's other assets, such
as the Dockerfile for the proxy's Docker image.


## metadata

The `metadata` section of the `.section-proxy.yaml` file is completely
section.io-specific. It's components are:

* `name`: Required. The name and version tag of the proxy as it would be referenced
in a section.io Application's `section.config.json` file.

* `image`: Required. The Registry-qualified Docker image that all containers in the
proxy Pod will be based upon. Containers within the Pod can have different entrypoints
but must all use the same image.

* `httpContainer`: Required. The name of the container in the Pod that will proxy
HTTP requests on its TCP port 80.

* `memoryContainer`: Required if more than one container is specified. The name
of the container in the Pod that can be allocated additional memory by the
section.io Operations team to better suit the needs of the website.

* `configurationMountPath`: Required. The directory within each container where the
proxy-specific subdirectory of section.io Application git repository will be mounted.
It will be mounted read-only.

* `logs.handler`: Required. The label applied to all the proxy's logs to identify
which code path be responsible for log parsing during ingestion.

* `logs.http.container`: Required. The name of the container that will write HTTP
access logs to one of its standard streams.

* `logs.http.stream`: Required. Either `stdout` or `stderr`. The standard stream
that will contain the HTTP access logs.

* `logs.additional`: Optional. An array of objects with mandatory `name`, `container`,
and `stream` members describing additional logs produced by the proxy. E.g. an
error.log, statistics, or detailed audit trail.

* `redeployTriggers.opsVariables`: Optional. An array of strings matching the
name of system operational variables that the pod cannot change on the fly.
Changing a matching variable value, will cause the pod to be (gracefully)
redeployed.

* `redeployTriggers.onAnyModuleFileChanges`: Optional. A boolean value that when
true indicates that any changes in the module's configuration will cause the pod
to be (gracefully) redeployed.

## spec

The `spec` section of the `.section-proxy.yaml` file is a constrained subset of
the [Kubernetes PodSpec](https://v1-9.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.9/#podspec-v1-core).
We only allow members that do not pose a risk to the platform, e.g. `serviceAccountName`
is not permitted since it could allow a proxy to elevate privileges.

We do not add section.io-custom members to this schema. If we need a section.io-specific
annotation, it should go in the `metadata` section described above and reference
the corresponding portion of `spec`, e.g. the `logs.http.container` member in `metadata`.

The currently permitted PodSpec subset includes:

* `volumes`: Optional. An array of Volume objects with only the `name` member permitted.
All volumes will be of type `emptyDir`.

* `containers`: Required. An array of at least one container. Each container must
have a `name`, and `resources`. Optionally containers may also define `volumeMounts`, `env`, and
`command`.

The `resources` expects both `requests` and `limits` members, which in turn expect
both `cpu` and `memory` to be specified. See the Kubernetes documentation for the
[meaning of CPU](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#meaning-of-cpu)
and the [meaning of memory](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#meaning-of-memory).

For the container referenced by the `memoryContainer` the `resources` define the
default values to use when not overridden by section.io Operations.

The `readinessProbe` and `livenessProve` members only accepts `initialDelaySeconds`.
If zero or omitted a default value will be used (20 seconds and 30 seconds respectively).
These settings are only honoured for the `httpContainer`.

The `volumeMounts` only accept `name` and `mountPath`. Volumes will be mounted read-write.

The `env` member requires `name` and `value`. The `valueFrom` member is not permitted.

## hpa

The optional `hpa` portion of the section proxy yaml is used to set defaults on the horizontal pod autoscaler. It tries to match the HPA v2beta1 API definition as closely as possible. We recommend leaving this section as default as the Section team will fine tune this. 

`apiVersion` - this is the API version the HPA spec is trying to match, only v2beta1 currently supported.
`spec` - contains the v2beta1 HPA spec https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#horizontalpodautoscalerspec-v2beta1-autoscaling

`spec.scaleTargetRef` - unsuported by section-proxy.yaml, Environment Provisioner controls this.
