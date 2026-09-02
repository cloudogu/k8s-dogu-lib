# Dogu format (v3beta1)

> **Beta**: `k8s.cloudogu.com/v3beta1` is a preview of the next Dogu-CR generation. It is served
> alongside `k8s.cloudogu.com/v2` (the current storage version) and the API server converts
> transparently between the two, so existing `v2` Dogu-CRs keep working unchanged. Setting
> `doguApiVersion: v3` opts a dogu into the upcoming Helm-based deployment strategy, which is not
> yet fully supported by the dogu operator. Until that lands, create Dogu-CRs as documented in
> [Dogu format](dogu_format_en.md) (`v2`) for production use.

All fields of a `v3beta1` Dogu-CR are described below and illustrated with examples. Only the
fields that are new or behave differently compared to `v2` are covered in detail; for everything
else see [Dogu format](dogu_format_en.md).

## Complete example

```yaml
apiVersion: k8s.cloudogu.com/v3beta1
kind: Dogu
metadata:
  name: usermgt
spec:
  name: usermgt
  doguNamespace: official
  version: 1.20.0-5
  doguApiVersion: v2
  mappedValues:
    logLevel: "debug"
  additionalIngressAnnotations: # deprecated
    nginx.ingress.kubernetes.io/proxy-body-size: "0"
  additionalMounts: # deprecated
    - sourceType: ConfigMap
      name: my-configmap
      volume: importHistory
      subfolder: "my-configmap-subfolder"
  exportMode: false
  resources: # deprecated
    minDataVolumeSize: 2Gi
  security: # deprecated
    appArmorProfile:
      localhostProfile: "localhost-profile"
      type: Localhost
    capabilities:
      add:
        - CAP_AUDIT_CONTROL
      drop:
        - CAP_SETGID
    readOnlyRootFileSystem: false
    runAsNonRoot: false
    seLinuxOptions:
      level: internal
      role: user_r
      type: user_t
      user: user_u
    seccompProfile:
      type: RuntimeDefault
  stopped: false
  pauseReconciliation: false
  supportMode: false # deprecated
  upgradeConfig:
    allowNamespaceSwitch: false
    forceUpgrade: false
  skipSchemaValidation: false
```

## Name

* Required
* Data type: string
* Content: Specifies the name of the Dogu without its namespace. In `v2` this was combined with
  the namespace into a single `namespace/name` string; in `v3beta1` the two are separate fields.
* Example: `"name": "usermgt"`

## DoguNamespace

* Required
* Data type: string
* Content: Specifies the namespace of the Dogu (the part before the slash in `v2`'s `name` field,
  e.g. `official`).
* Example: `"doguNamespace": "official"`

## Version

* Required
* Data type: string
* Content: Specifies the version of the Dogu.
* Example: `"version": "1.20.0-5"`

## DoguApiVersion

* Optional
* Data type: Enum <v2; v3>
* Default: `v2`
* Content: Tells the dogu operator's reconciler which internal deployment routine to run for this
  dogu.
    - `v2` - builds Kubernetes resources ad-hoc from the dogu's attributes, the same way
      `k8s.cloudogu.com/v2` dogus are deployed today. `resources`, `security`, `supportMode` and
      `additionalMounts` apply.
    - `v3` - **not yet supported by the dogu operator.** Once available, this deploys and manages
      the dogu as a Helm release, configured via `values` instead of `resources`/`security`/
      `additionalMounts`.
* Example: `"doguApiVersion": "v2"`

## Values

* Optional
* Data type: Object
* Content: Helm `values.yaml` overrides for this dogu's chart. Only meaningful when
  `doguApiVersion` is `v3`; ignored otherwise. The structure is arbitrary and validated by the
  chart itself, not by this schema.
* Example:

```
values:
  replicaCount: 2
```

## MappedValues

* Optional
* Data type: Object (string-to-string map)
* Content: MappedValues maps metadata values to specific values from the helm chart. This is
  typically used to change all container log levels by setting one single value.
* Example:

```
mappedValues:
  logLevel: "debug"
```

## AdditionalIngressAnnotations

This field is `deprecated`. It is only honored for `doguApiVersion: v2`.

* Optional
* Data type: string
* Content: AdditionalIngressAnnotations provides additional annotations that get included into the dogu's ingress rules.
* Example:

```
additionalIngressAnnotations:
  nginx.ingress.kubernetes.io/proxy-body-size: "0"
```

## AdditionalMounts

This field is `deprecated`. It is only honored for `doguApiVersion: v2`; Helm-deployed dogus
(`doguApiVersion: v3`) configure additional mounts via `values` instead.

* Optional
* Data type: Array<DataMount>
* Content: AdditionalMounts provides the possibility to mount additional data into the dogu.
* Example:

```
  additionalMounts:
    - sourceType: ConfigMap
      name: my-configmap
      volume: importHistory
      subfolder: "my-configmap-subfolder"
    - sourceType: Secret
      name: my-secret
      volume: importHistory
```

### DataMount

A DataMount can contain the following fields:

#### SourceType

* Required
* Data type: Enum <ConfigMap; Secret>
* Content: SourceType defines where the data is coming from.
  Valid options are:
    - ConfigMap - data stored in a kubernetes ConfigMap.
    - Secret - data stored in a kubernetes Secret.
* Example: `"sourceType": ConfigMap`

#### Name

* Required
* Data type: String
* Content: Name is the name of the data source.
* Example: `"name": my-configmap`

#### Volume

* Required
* Data type: String
* Content: Volume is the name of the volume to which the data should be mounted. It is defined in the respective
  dogu.json.
* Example: `"volume": importHistory`

#### Subfolder

* Optional
* Data type: String
* Content: Subfolder defines a subfolder in which the data should be put within the volume.
* Example: `"subfolder": "my-configmap-subfolder"`

## ExportMode

* Optional
* Data type: String
* Content: ExportMode indicates whether the dogu should be in "export mode". If true, the operator will spawn an
  exporter sidecar container along with a new volume mount to aid the migration process from one Cloudogu EcoSystem to
  another.
* Example: `"exportMode": false`

## Resources

This field is `deprecated`. It is only honored for `doguApiVersion: v2`; Helm-deployed dogus
(`doguApiVersion: v3`) configure volume sizing via `values` instead.

* Optional
* Data type: Object
* Content: Resources of the dogu (e.g. minDataVolumeSize)
* Example:

```
resources:
  minDataVolumeSize: 2Gi
```

### MinDataVolumeSize

* Optional
* Data type: String
* Content: MinDataVolumeSize represents the desired minimum size of the volume. Increasing this value may lead to an automatic volume
  expansion. This includes a downtime for the respective dogu. The default size for volumes is "2Gi".
  It is not possible to lower the volume size after an expansion. This will introduce an inconsistent state for the
  dogu.
* Example: `"minDataVolumeSize": 2Gi`

## Security

This field is `deprecated`. It is only honored for `doguApiVersion: v2`; Helm-deployed dogus
(`doguApiVersion: v3`) configure container security via `values` instead.

* Optional
* Data type: Object
* Content: Security overrides security policies defined in the dogu descriptor. These fields can be used to further
  reduce a dogu's attack surface.

Security can contain the following attributes:

### AppArmorProfile

* Optional
* Data type: Object
* Content: AppArmorProfile is the AppArmor options to use by this container.
* Example:

```
appArmorProfile:
  localhostProfile: "localhost-profile"
  type: Localhost
```

#### Type

* Required
* Data type: Enum <Localhost; RuntimeDefault; Unconfined>
* Content: Type indicates which kind of AppArmor profile will be applied.
  Valid options are:
    - Localhost - a profile pre-loaded on the node.
    - RuntimeDefault - the container runtime's default profile.
    - Unconfined - no AppArmor enforcement.

#### LocalhostProfile

* Optional
* Data type: String
* Content: LocalhostProfile indicates a profile loaded on the node that should be used.
  The profile must be preconfigured on the node to work. Must match the loaded name of the profile.
  Must be set if and only if type is "Localhost".

### Capabilities

* Optional
* Data type: Object
* Content: Capabilities sets the allowed and dropped capabilities for the dogu. The dogu should not use more than the
  configured capabilities here, otherwise failure may occur at start-up or at run-time. Each capability represents a
  POSIX capabilities type. See docs at https://manned.org/capabilities.7
* Example:

```
capabilities:
  add:
    - CAP_AUDIT_CONTROL
  drop:
    - CAP_SETGID
```

#### Add

* Optional
* Data type: Array<String>
* Content: Add contains the capabilities that should be allowed to be used in a container. This list is optional.

#### Drop

* Optional
* Data type: Array<String>
* Content: Drop contains the capabilities that should be blocked from being used in a container. This list is optional.

### ReadOnlyRootFileSystem

* Optional
* Data type: boolean
* Content: ReadOnlyRootFileSystem mounts the container's root filesystem as read-only. The dogu must support accessing
  the root file system by only reading otherwise the dogu start may fail. This flag is optional and defaults to nil.
  If nil, the value defined in the dogu descriptor is used.
* Example: `"readOnlyRootFileSystem": true`

### RunAsNonRoot

* Optional
* Data type: boolean
* Content: RunAsNonRoot indicates that the container must run as a non-root user. The dogu must support running as
  non-root user otherwise the dogu start may fail. This flag is optional and defaults to nil.
  If nil, the value defined in the dogu descriptor is used.
* Example: `"runsNonRoot": true`

### SeLinuxOptions

* Optional
* Data type: Object
* Content: SELinuxOptions is the SELinux context to be applied to the container.
  If unspecified, the container runtime will allocate a random SELinux context for each container, which is kubernetes
  default behaviour.
* Example:

```
seLinuxOptions:
  level: internal
  role: user_r
  type: user_t
  user: user_u
```

#### Level

* Optional
* Data type: string
* Content: Level is SELinux level label that applies to the container.

#### Role

* Optional
* Data type: string
* Content: Role is a SELinux role label that applies to the container.

#### Type

* Optional
* Data type: string
* Content: Type is a SELinux type label that applies to the container.

#### User

* Optional
* Data type: string
* Content: User is a SELinux user label that applies to the container.

### SeccompProfile

* Optional
* Data type: Object
* Content: SeccompProfile is the seccomp options to use by this container.
* Example:

```
seccompProfile:
  localhostProfile: "localhost-profile"
  type: Localhost
```

#### Type

* Required
* Data type: Enum <Localhost; RuntimeDefault; Unconfined>
* Content: Type indicates which kind of AppArmor profile will be applied.
  Valid options are:
    - Localhost - a profile defined in a file on the node should be used.
    - RuntimeDefault - the container runtime default profile should be used.
    - Unconfined - no profile should be applied.

#### LocalhostProfile

* Optional
* Data type: string
* Content: LocalhostProfile indicates a profile defined in a file on the node should be used.
  The profile must be preconfigured on the node to work.
  Must be a descending path, relative to the kubelet's configured seccomp profile location.
  Must be set if type is "Localhost". Must NOT be set for any other type.

## Stopped

* Optional
* Data type: boolean
* Content: Stopped indicates whether the dogu should be running (stopped=false) or not (stopped=true).
* Example: `"stopped": true`

## PauseReconciliation

* Optional
* Data type: boolean
* Content: PauseReconciliation indicates whether the reconciliation loop should be running
  (pauseReconciliation=false) or not (pauseReconciliation=true). The validation step always keeps
  running regardless of this setting.
* Example: `"pauseReconciliation": true`

## SupportMode

This field is `deprecated`. It is only honored for `doguApiVersion: v2`.

* Optional
* Data type: boolean
* Content: SupportMode indicates whether the dogu should be restarted in the support mode (f. e. to recover manually
  from a crash loop).
* Example: `"supportMode": true`

## UpgradeConfig

* Optional
* Data type: Object
* Content: UpgradeConfig contains options to manipulate the upgrade process.
* Example:

```
upgradeConfig:
  allowNamespaceSwitch: false
  forceUpgrade: false
```

### AllowNamespaceSwitch

* Optional
* Data type: boolean
* Content: AllowNamespaceSwitch lets a dogu switch its dogu namespace during an upgrade. The dogu must be technically
  the same dogu which did reside in a different namespace. The remote dogu's version must be equal to or greater than
  the version of the local dogu.

### ForceUpgrade

* Optional
* Data type: boolean
* Content: ForceUpgrade allows to install the same or even lower dogu version than already is installed. Please note,
  that possible data loss may occur by inappropriate dogu downgrading.

## SkipSchemaValidation

* Optional
* Data type: boolean
* Content: SkipSchemaValidation indicates whether the dogu's schema validation should be skipped.
* Example: `"skipSchemaValidation": true`

## Status

The `status` block of a `v3beta1` Dogu-CR is written exclusively by the dogu operator and reflects
the operator's observed state of the dogu.

### Example

```yaml
status:
  installedVersion: 1.20.0-5
  appVersion: 4.2.1
  stopped: false
  exportMode: false
  startedAt: "2026-08-01T10:00:00Z"
  conditions:
    - type: ready
      status: "True"
    - type: healthy
      status: "True"
    - type: stopped
      status: "False"
      reason: NotStopped
    - type: pauseReconciliation
      status: "False"
      reason: ReconciliationEnabled
  volumeStatus:
    - claimName: usermgt-data
      size: 2Gi
      desiredSize: 2Gi
      state: ready
```

### InstalledVersion

* Data type: string
* Content: The currently installed version of the dogu chart (e.g. `1.20.0-5`).
* Example: `"installedVersion": "1.20.0-5"`

### AppVersion

* Data type: string
* Content: The version of the application contained in the dogu (e.g. `4.2.1`). Unlike
  `installedVersion`, which reports the dogu chart's version.
* Example: `"appVersion": "4.2.1"`

### StartedAt

* Data type: timestamp
* Content: The time of the dogu's last (re)start.

### Stopped

* Data type: boolean
* Content: Reports the observed state corresponding to `spec.stopped`. May briefly differ from the
  `spec` field while a stop or start operation is in progress.

### ExportMode

* Data type: boolean
* Content: Indicates whether the dogu's export mode is currently actually active. Observed state
  corresponding to `spec.exportMode`.

### Conditions

* Data type: Array<Condition> (standard Kubernetes `metav1.Condition`)
* Content: Conditions represent the machine-readable state of the dogu. Each condition has a
  `type`, a `status` (`True`/`False`/`Unknown`), a `reason`, and optionally a `message` and a
  `lastTransitionTime`. `v3beta1` uses the following condition types:
    - `ready` - `True` once no asynchronous operation is currently running for the dogu.
      `reason` while `False`: `Installing`, `Upgrading`, `ResizingPVC`, `ChangingExportMode`,
      `Starting`, `Stopping`, `Deleting`.
    - `healthy` - `True` when the workloads provided by the dogu are running and ready.
      `reason` while `False`: `Stopped`, `WorkloadsNotReady`.
    - `stopped` - reflects whether all of the dogu's workloads are actually stopped.
      `reason`: `NotStopped`, `PartiallyStopped`.
    - `pauseReconciliation` - reflects `spec.pauseReconciliation`.
      `reason`: `ReconciliationEnabled`, `ReconciliationPaused`.
    - `valid` - `False` when the Dogu-CR cannot be processed for functional reasons (independent
      of plain schema validity). `reason`: `SchemaInvalid`, `UnsupportedField`, `VolumeShrink`,
      `StorageClassImmutable`, `NoUpgradePath`, `ExportModeOnInstall`.
    - `chartAvailable` - `False` when the Helm chart required for `doguApiVersion: v3` cannot be
      loaded. `reason`: `ChartNotFound`, `DownloadFailed`, `Unauthorized`.
    - `updatePending` - `True` when there is a change to the Dogu-CR that is currently not being
      applied (e.g. because the dogu is stopped or reconciliation is paused).
      `reason`: `Stopped`, `PauseReconciliation`.
    - `schemaValidationSkipped` - reflects `spec.skipSchemaValidation`.
      `reason`: `SchemaValidationEnabled`, `SchemaValidationDisabled`.

  The `ready`, `healthy`, and `pauseReconciliation` conditions already existed in `v2` (some with
  different `reason` values); all others are new in `v3beta1`. `kubectl get dogu` shows `ready`,
  `healthy`, `stopped`, and `pauseReconciliation` as dedicated columns.

### VolumeStatus

* Data type: Array<VolumeStatus>
* Content: VolumeStatus contains the observed state for every PersistentVolumeClaim used by the
  dogu.
* Example:

```
volumeStatus:
  - claimName: usermgt-data
    size: 2Gi
    desiredSize: 2Gi
    state: ready
```

#### ClaimName

* Data type: string
* Content: Name of the PersistentVolumeClaim.

#### Size

* Data type: string (Quantity)
* Content: The volume's current actual size.

#### DesiredSize

* Data type: string (Quantity)
* Content: The target size for the volume desired by `spec`.

#### State

* Data type: Enum <ready>
* Content: State of the volume. Currently only the `ready` value exists; further states (e.g. for
  an in-progress resize) are planned for future versions.

#### MountedBy

* Data type: Array<Object>
* Content: References to the Kubernetes objects (e.g. Pods) currently mounting the volume.

### Deprecated fields

The following fields still exist for compatibility with `v2`, but are marked as `deprecated` and
are planned for removal in a future major version. They are identical in content to their
counterparts in [Dogu format](dogu_format_en.md) (`v2`) and are superseded by the fields and
conditions described above:

* `status` (string) - the old, textual overall status of the dogu. Superseded by the `ready`
  condition. Possible values: `""` (not installed), `installing`, `upgrading`, `deleting`,
  `installed`, `resizing PVC`, `starting`, `stopping`, `changing export-mode`,
  `change data mounts`.
* `requeueTime` - wait time until the next requeue of a running asynchronous operation.
* `requeuePhase` - current phase of a running asynchronous operation.
* `health` - the old health status of the dogu (`""`, `available`, `unavailable`, `unknown`).
  Superseded by the `healthy` condition.
* `dataVolumeSize` - the old size of the data volume. Superseded by `volumeStatus`.
