package v3beta1

import (
	"time"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Compile-Time Interface Assertions
var _ conditions.Getter = &Dogu{}
var _ conditions.Setter = &Dogu{}

// These constants are exported for use in other packages
// nolint:unused
//
//goland:noinspection GoUnusedConst
const (
	// DoguLabelName is used to select a dogu pod by name.
	DoguLabelName = "dogu.name"
	// DoguLabelVersion is used to select a dogu pod by version.
	DoguLabelVersion = "dogu.version"
)

// DoguApiVersion tells the operator's reconciler which internal deployment routine to
// run for a dogu.
// +enum
type DoguApiVersion string

// These constants are exported for use in other packages
// nolint:unused
//
//goland:noinspection GoUnusedConst
const (
	// DoguApiVersionV2 builds Kubernetes resources ad-hoc from the dogu's attributes,
	// the same way k8s.cloudogu.com/v2 dogus are deployed.
	DoguApiVersionV2 DoguApiVersion = "v2"
	// DoguApiVersionV3 deploys and manages the dogu as a Helm release, using
	// DoguSpec.Values as values.yaml overrides.
	DoguApiVersionV3 DoguApiVersion = "v3"
)

// DoguSpec defines the desired state of a Dogu
type DoguSpec struct {
	// Name of the dogu without the namespace (e.g., ldap)
	Name string `json:"name,omitempty"`
	// DoguNamespace of the dogu without the actual dogu name (e.g., official)
	DoguNamespace string `json:"doguNamespace,omitempty"`
	// Version of the dogu (e.g. 2.4.48-3)
	Version string `json:"version,omitempty"`

	// DoguApiVersion tells the operator's reconciler which internal deployment routine to run for this
	// dogu. "v2" builds Kubernetes resources ad-hoc from the dogu's attributes (compatible with
	// k8s.cloudogu.com/v2 dogus). "v3" deploys the dogu as a Helm release, configured via Values.
	// +kubebuilder:validation:Enum=v2;v3
	// +kubebuilder:default=v2
	// +optional
	DoguApiVersion DoguApiVersion `json:"doguApiVersion,omitempty"`
	// Values contains Helm values.yaml overrides for this dogu's chart. Only meaningful when
	// DoguApiVersion is "v3"; ignored otherwise. The structure is arbitrary and validated by the
	// chart itself, not by this schema.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values apiextensionsv1.JSON `json:"values,omitempty"`
	// MappedValues maps metadata values to specific values from the helm chart.
	// This is typically used to change all container log levels by settings one single value.
	MappedValues map[string]string `json:"mappedValues,omitempty"`
	// Deprecated: Resources is only honored by the "legacy" deployment strategy. Helm-deployed
	// dogus (DoguApiVersion=helm) configure volume sizing via Values instead.
	Resources DoguResources `json:"resources,omitempty"`
	// Deprecated: Security is only honored by the "legacy" deployment strategy. Helm-deployed
	// dogus (DoguApiVersion=helm) configure container security via Values instead.
	// +optional
	Security Security `json:"security,omitempty"`
	// Deprecated: SupportMode is only honored by the "legacy" deployment strategy.
	SupportMode bool `json:"supportMode,omitempty"`
	// ExportMode indicates whether the dogu should be in "export mode". If true, the operator will spawn an exporter sidecar
	// container along with a new volume mount to aid the migration process from one Cloudogu EcoSystem to another.
	ExportMode bool `json:"exportMode,omitempty"`
	// Stopped indicates whether the dogu should be running (stopped=false) or not (stopped=true).
	Stopped bool `json:"stopped,omitempty"`
	// PauseReconciliation indicates whether the reconciliation loop should be running (pauseReconciliation=false) or not (pauseReconciliation=true).
	// The validation step should always be running.
	PauseReconciliation bool `json:"pauseReconciliation,omitempty"`
	// UpgradeConfig contains options to manipulate the upgrade process.
	UpgradeConfig UpgradeConfig `json:"upgradeConfig,omitempty"`
	// Deprecated: This field is not used anymore. It was used to provide additional configuration for the nginx ingress controller.
	// The traefik ingress controller doesn't need this configuration.
	// AdditionalIngressAnnotations provides additional annotations that get included into the dogu's ingress rules.
	AdditionalIngressAnnotations IngressAnnotations `json:"additionalIngressAnnotations,omitempty"`
	// Deprecated: AdditionalMounts is only honored by the "legacy" deployment strategy. Helm-deployed
	// dogus (DoguApiVersion=helm) configure additional mounts via Values instead.
	// +optional
	AdditionalMounts []DataMount `json:"additionalMounts,omitempty" patchStrategy:"replace"` // no unique identifier, so we can't use merge
	// SkipSchemaValidation indicates whether the schema validation should be skipped.
	SkipSchemaValidation bool `json:"skipSchemaValidation,omitempty"`
}

// DataSourceType defines the supported source types of additional data mounts.
// +enum
type DataSourceType string

// These constants are exported for use in other packages
// nolint:unused
//
//goland:noinspection GoUnusedConst
const (
	// DataSourceConfigMap mounts a config map as a data source.
	DataSourceConfigMap DataSourceType = "ConfigMap"
	// DataSourceSecret mounts a secret as a data source.
	DataSourceSecret DataSourceType = "Secret"
)

// DataMount is a description of what data should be mounted to a specific Dogu volume (already defined in dogu.json).
type DataMount struct {
	// SourceType defines where the data is coming from.
	// Valid options are:
	//   ConfigMap - data stored in a kubernetes ConfigMap.
	//   Secret - data stored in a kubernetes Secret.
	// +kubebuilder:validation:Enum=ConfigMap;Secret
	SourceType DataSourceType `json:"sourceType"`
	// Name is the name of the data source.
	Name string `json:"name"`
	// Volume is the name of the volume to which the data should be mounted. It is defined in the respective dogu.json.
	Volume string `json:"volume"`
	// Subfolder defines a subfolder in which the data should be put within the volume.
	// +optional
	Subfolder string `json:"subfolder,omitempty"`
}

// IngressAnnotations are annotations of nginx-ingress rules.
type IngressAnnotations map[string]string

// UpgradeConfig contains configuration hints for the dogu operator regarding aspects during the upgrade of dogus.
type UpgradeConfig struct {
	// AllowNamespaceSwitch lets a dogu switch its dogu namespace during an upgrade. The dogu must be technically the
	// same dogu which did reside in a different namespace. The remote dogu's version must be equal to or greater than
	// the version of the local dogu.
	AllowNamespaceSwitch bool `json:"allowNamespaceSwitch,omitempty"`
	// ForceUpgrade allows to install the same or even lower dogu version than already is installed. Please note, that
	// possible data loss may occur by inappropriate dogu downgrading.
	ForceUpgrade bool `json:"forceUpgrade,omitempty"`
}

// DoguResources defines the physical resources used by the dogu.
//
// Deprecated: Resources is only honored by the "legacy" deployment strategy.
// +kubebuilder:validation:XValidation:rule="(!has(oldSelf.storageClassName) && !has(self.storageClassName)) || (has(oldSelf.storageClassName) && has(self.storageClassName))", message="StorageClassName cannot be set or unset after creation"
type DoguResources struct {
	// DataVolumeSize represents the desired size of the volume. Increasing this value leads to an automatic volume
	// expansion. This includes a downtime for the respective dogu. The default size for volumes is "2Gi".
	// Attempts to lower the size of an existing Dogu will be ignored.
	// Has the format of a resource.Quantity.
	//
	// Deprecated. Now acts the same as MinDataVolumeSize and will soon be replaced by it.
	// It is recommended to not write this field and read the value by calling Dogu.GetMinDataVolumeSize which will consider MinDataVolumeSize as well.
	// If both this and MinDataVolumeSize are set, MinDataVolumeSize takes precedent.
	DataVolumeSize string `json:"dataVolumeSize,omitempty"`
	// MinDataVolumeSize represents the minimum desired size of the volume. Increasing this value leads to an automatic volume
	// expansion. This includes a downtime for the respective dogu. The default size for volumes is "2Gi".
	// Attempts to lower the size of an existing Dogu will be ignored.
	//
	// The value of MinDataVolumeSize takes precedent over DataVolumeSize.
	// To consider both values when reading, call Dogu.GetMinDataVolumeSize.
	MinDataVolumeSize resource.Quantity `json:"minDataVolumeSize,omitempty"`
	// StorageClassName specifies the storage class to be used for the data volume.
	// For the difference between null and empty string, see the appropriate kubernetes documentation:
	// https://kubernetes.io/docs/concepts/storage/persistent-volumes/#class-1
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="StorageClassName is immutable"
	// +kubebuilder:validation:MaxLength=253
	StorageClassName *string `json:"storageClassName,omitempty"`
}

type HealthStatus string

// These constants are exported for use in other packages
// nolint:unused
//
//goland:noinspection GoUnusedConst
const (
	PendingHealthStatus     HealthStatus = ""
	AvailableHealthStatus   HealthStatus = "available"
	UnavailableHealthStatus HealthStatus = "unavailable"
	UnknownHealthStatus     HealthStatus = "unknown"
)

type VolumeStatusState string

const (
	VolumeStatusStateReady VolumeStatusState = "ready"
	// TODO More states? Resizing?
)

type MountRef struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type VolumeStatus struct {
	ClaimName   string            `json:"claimName"`
	Size        resource.Quantity `json:"size,omitempty"`
	DesiredSize resource.Quantity `json:"desiredSize,omitempty"`
	State       VolumeStatusState `json:"state,omitempty"`
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=kind
	// +listMapKey=name
	MountedBy []corev1.TypedLocalObjectReference `json:"mountedBy,omitempty"`
}

// DoguStatus defines the observed state of a Dogu.
type DoguStatus struct {
	// Status represents the state of the Dogu in the ecosystem
	// Deprecated, should be removed at the next major update
	Status string `json:"status"`
	// RequeueTime contains time necessary to perform the next requeue
	// Deprecated, should be removed at the next major update
	RequeueTime time.Duration `json:"requeueTime"`
	// RequeuePhase is the actual phase of the dogu resource used for a currently running async process.
	// Deprecated, should be removed at next major update
	RequeuePhase string `json:"requeuePhase"`
	// StartedAt contains the time of the last restart of the dogu.
	StartedAt metav1.Time `json:"startedAt,omitempty"`
	// Health describes the health status of the dogu
	// Deprecated, should be removed at next major update
	Health HealthStatus `json:"health,omitempty"`
	// InstalledVersion of the dogu chart (e.g. 2.4.48)
	InstalledVersion string `json:"installedVersion,omitempty"`
	// AppVersion of the dogu chart (e.g. 4.2.1)
	AppVersion string `json:"appVersion,omitempty"`
	// Stopped shows if the dogu has been stopped or not.
	Stopped bool `json:"stopped,omitempty"`
	// ExportMode shows if the export mode of the dogu is currently active.
	ExportMode bool `json:"exportMode,omitempty"`
	// DataVolumeSize shows the current size of the mounted data volume
	// Deprecated, replaced by volumes and should be removed at the next major update.
	DataVolumeSize *resource.Quantity `json:"dataVolumeSize,omitempty"`
	// a list of conditions TRUE|FALSE
	// e.g. MeetsMinimumDataVolumeSize -> True if status.dataVolumeSize >= spec.minDataVolumeSize
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	// VolumesStatus contains relevant information about all persistent volume claims used by the dogu.
	// +patchMergeKey=claimName
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=claimName
	VolumeStatus []VolumeStatus `json:"volumeStatus,omitempty"`
}

// These constants are exported for use in other packages
// nolint:unused
//
//goland:noinspection GoUnusedConst
const (
	DoguStatusNotInstalled       = ""
	DoguStatusInstalling         = "installing"
	DoguStatusUpgrading          = "upgrading"
	DoguStatusDeleting           = "deleting"
	DoguStatusInstalled          = "installed"
	DoguStatusPVCResizing        = "resizing PVC"
	DoguStatusStarting           = "starting"
	DoguStatusStopping           = "stopping"
	DoguStatusChangingExportMode = "changing export-mode"
	DoguStatusChangingDataMounts = "change data mounts"
)

// Conditions present in v2 and v3beta1
const (
	// ConditionReady
	// Reasons v3beta1: Installing, Upgrading, ResizingPVC, ChangingExportMode, Starting, Stopping, Deleting
	// Reasons v2: ReconcileSuccess, ReconcileFail, HasToReconcile
	ConditionReady = "ready"
	// ConditionHealthy
	// Reasons v3beta1: Stopped, WorkloadsNotReady
	// Reasons v2: StoppingOperator, Deleting, DoguIsNotHealthy, DoguIsHealthy, Upgrading
	ConditionHealthy             = "healthy"             // Needs to be translated
	ConditionPauseReconciliation = "pauseReconciliation" // Same condition with equal reasons
)

// New conditions
const (
	ConditionStopped                 = "stopped"
	ConditionPausedValid             = "valid"
	ConditionChartAvailable          = "chartAvailable"
	ConditionUpdatePending           = "updatePending"
	ConditionSchemaValidationSkipped = "schemaValidationSkipped"
	ConditionExportModeActive        = "exportModeActive"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Spec Version",type="string",JSONPath=".spec.version",description="The desired version of the dogu chart"
// +kubebuilder:printcolumn:name="Installed Version",type="string",JSONPath=".status.installedVersion",description="The current version of the dogu chart"
// +kubebuilder:printcolumn:name="App Version",type="string",JSONPath=".status.appVersion",description="The current version of the dogu application"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='ready')].status",description="Whether the resource is ready in the current state"
// +kubebuilder:printcolumn:name="Healthy",type="string",JSONPath=".status.conditions[?(@.type=='healthy')].status",description="Whether the resource is healthy in the current state"
// +kubebuilder:printcolumn:name="ApiVersion",type="string",JSONPath=".spec.doguApiVersion",description="The api version for this dogu"
// +kubebuilder:printcolumn:name="Stopped",type="string",JSONPath=".status.conditions[?(@.type=='stopped')].status",description="Whether the dogu is stopped or not"
// +kubebuilder:printcolumn:name="Pause Reconciliation",type="string",JSONPath=".status.conditions[?(@.type=='pauseReconciliation')].status",description="Whether the resource is ready in the current state"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of the resource"

// Dogu is the Schema for the dogus API
type Dogu struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DoguSpec   `json:"spec,omitempty"`
	Status DoguStatus `json:"status,omitempty"`
}

func (d *Dogu) GetConditions() []metav1.Condition {
	return d.Status.Conditions
}

func (d *Dogu) SetConditions(c []metav1.Condition) {
	d.Status.Conditions = c
}

// GetSimpleDoguName returns the name of the dogu as a dogu.SimpleName.
func (d *Dogu) GetSimpleDoguName() cescommons.SimpleName {
	return cescommons.SimpleName(d.Spec.Name)
}

// GetSimpleNameVersion returns the dogu's dogu.SimpleNameVersion.
func (d *Dogu) GetSimpleNameVersion() (cescommons.SimpleNameVersion, error) {
	version, err := core.ParseVersion(d.Spec.Version)
	if err != nil {
		return cescommons.SimpleNameVersion{}, err
	}

	return cescommons.NewSimpleNameVersion(d.GetSimpleDoguName(), version), nil
}

// +kubebuilder:object:root=true

// DoguList contains a list of Dogu
type DoguList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Dogu `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Dogu{}, &DoguList{})
}
