# Dogu format

The Dogu-CR can be used to install Cloudogu-Dogus in a Kubernetes cluster with the dogu operator.
Various settings can be made that are evaluated in addition to dogu.json and the dogu configuration.

All fields of a Dogu-CR are described below and illustrated with examples.

## Komplettes Beispiel

```yaml
apiVersion: k8s.cloudogu.com/v2
kind: Dogu
metadata:
  name: postfix
spec:
  name: official/usermgt
  version: 1.20.0-5
  additionalIngressAnnotations:
    nginx.ingress.kubernetes.io/proxy-body-size: "0"
  additionalMounts:
    - sourceType: ConfigMap
      name: my-configmap
      volume: importHistory
      subfolder: "my-configmap-subfolder"
  exportMode: false
  resources:
    minDataVolumeSize: 2Gi
  security:
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
  supportMode: false
  upgradeConfig:
    allowNamespaceSwitch: false
    forceUpgrade: false
```

## Name

* Required
* Data type: string
* Content: Specifies the name including the namespace of the Dogu.
* Example: `"name": "official/usermgt"`

## Version

* Required
* Data type: string
* Content: Specifies the version of the Dogu.
* Example: `"version": "1.20.0-5"`

## AdditionalIngressAnnotations

* Optional
* Data type: string
* Content: AdditionalIngressAnnotations provides additional annotations that get included into the dogu's ingress rules.
* Example:

```
additionalIngressAnnotations:
  nginx.ingress.kubernetes.io/proxy-body-size: "0"
```

## AdditionalMounts

* Optional
* Data type: Array<DataMount>
* Content: Data provides the possibility to mount additional data into the dogu.
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

## SupportMode

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
  the
  same dogu which did reside in a different namespace. The remote dogu's version must be equal to or greater than
  the version of the local dogu.

### ForceUpgrade

* Optional
* Data type: boolean
* Content: ForceUpgrade allows to install the same or even lower dogu version than already is installed. Please note,
  that possible data loss may occur by inappropriate dogu downgrading.
