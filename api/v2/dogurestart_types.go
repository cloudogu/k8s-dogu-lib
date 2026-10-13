/*
This file was generated with "make generate-deepcopy".
*/
package v2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

var _ conditions.Getter = &DoguRestart{}
var _ conditions.Setter = &DoguRestart{}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DoguRestartSpec defines the desired state of DoguRestart
type DoguRestartSpec struct {
	// DoguName references the dogu that should get restarted.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Dogu name is immutable"
	DoguName string `json:"doguName"`
}

// DoguRestartStatus defines the observed state of DoguRestart
type DoguRestartStatus struct {
	// Phase tracks the state of the restart process.
	// It is only used for restarts of v2 dogus. Restarts of v3 dogus use Conditions instead.
	Phase RestartStatusPhase `json:"phase,omitempty"`
	// Conditions track the state of the restart process of v3 dogus.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

type RestartStatusPhase string

func (rsp RestartStatusPhase) IsFailed() bool {
	return rsp != RestartStatusPhaseNew && !rsp.isInProgress() && !rsp.isSuccessful()
}

func (rsp RestartStatusPhase) isInProgress() bool {
	return rsp == RestartStatusPhaseStarting || rsp == RestartStatusPhaseStopping || rsp == RestartStatusPhaseStopped
}

func (rsp RestartStatusPhase) isSuccessful() bool {
	return rsp == RestartStatusPhaseCompleted
}

const (
	RestartStatusPhaseNew           RestartStatusPhase = ""
	RestartStatusPhaseStopping      RestartStatusPhase = "stopping"
	RestartStatusPhaseStopped       RestartStatusPhase = "stopped"
	RestartStatusPhaseStarting      RestartStatusPhase = "starting"
	RestartStatusPhaseCompleted     RestartStatusPhase = "completed"
	RestartStatusPhaseDoguNotFound  RestartStatusPhase = "dogu not found"
	RestartStatusPhaseFailedGetDogu RestartStatusPhase = "failed getting dogu"
	RestartStatusPhaseFailedStop    RestartStatusPhase = "stop failed"
	RestartStatusPhaseFailedStart   RestartStatusPhase = "start failed"
)

// Conditions for restarts of v3 dogus.
const (
	// DoguRestartConditionRestarted is true once the restart annotation is set on all workloads of the dogu.
	DoguRestartConditionRestarted = "Restarted"
	// DoguRestartConditionReady is true once all workloads of the dogu are rolled out after the restart.
	DoguRestartConditionReady = "Ready"
)

// Reasons for the DoguRestart conditions.
const (
	// DoguRestartReasonSucceeded is the positive reason for all DoguRestart conditions.
	DoguRestartReasonSucceeded = "Succeeded"
	// DoguRestartReasonSuperseded is used if a newer restart of the same dogu was already applied.
	DoguRestartReasonSuperseded = "Superseded"
	// DoguRestartReasonDoguNotFound is used if the dogu referenced in spec.doguName does not exist.
	DoguRestartReasonDoguNotFound = "DoguNotFound"
	// DoguRestartReasonDoguV3Disabled is used if the dogu is a v3 dogu but v3 support is disabled in the operator.
	DoguRestartReasonDoguV3Disabled = "DoguV3Disabled"
	// DoguRestartReasonDoguNotIdle is used while the dogu is installing or upgrading and the restart waits for it.
	DoguRestartReasonDoguNotIdle = "DoguNotIdle"
	// DoguRestartReasonPatchFailed is used if the restart annotation could not be set on a workload.
	DoguRestartReasonPatchFailed = "PatchFailed"
	// DoguRestartReasonRolloutInProgress is used while the workloads are rolled out after the restart.
	DoguRestartReasonRolloutInProgress = "RolloutInProgress"
	// DoguRestartReasonDoguStopped is used if the dogu is stopped. The restart is applied but no pods are started.
	DoguRestartReasonDoguStopped = "DoguStopped"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName="dr"
// +kubebuilder:printcolumn:name="Dogu",type="string",JSONPath=".spec.doguName",description="The name of the dogu"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="The current phase of the dogu restart"
// +kubebuilder:printcolumn:name="Restarted",type="string",JSONPath=".status.conditions[?(@.type=='Restarted')].status",description="Whether the restart is applied to all workloads"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status",description="Whether all workloads are rolled out after the restart"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of the resource"

// DoguRestart is the Schema for the dogurestarts API
type DoguRestart struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DoguRestartSpec   `json:"spec,omitempty"`
	Status DoguRestartStatus `json:"status,omitempty"`
}

func (dr *DoguRestart) GetConditions() []metav1.Condition {
	return dr.Status.Conditions
}

func (dr *DoguRestart) SetConditions(c []metav1.Condition) {
	dr.Status.Conditions = c
}

// +kubebuilder:object:root=true

// DoguRestartList contains a list of DoguRestart
type DoguRestartList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DoguRestart `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DoguRestart{}, &DoguRestartList{})
}
