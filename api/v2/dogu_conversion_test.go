package v2

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/cloudogu/cesapp-lib/core"
	v3beta1 "github.com/cloudogu/k8s-dogu-lib/v2/api/v3beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func quantityPtr(s string) *resource.Quantity {
	return new(resource.MustParse(s))
}

func newFullV2TestDogu() *Dogu {
	return &Dogu{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "postgresql",
			Namespace:   "ecosystem",
			Annotations: map[string]string{"custom.example.com/note": "unrelated annotation"},
		},
		Spec: DoguSpec{
			Name:    "official/postgresql",
			Version: "1.2.3-4",
			Resources: DoguResources{
				DataVolumeSize:    "5Gi",
				MinDataVolumeSize: resource.MustParse("5Gi"),
				StorageClassName:  new("longhorn"),
			},
			Security: Security{
				Capabilities: Capabilities{
					Add:  []core.Capability{core.AuditControl},
					Drop: []core.Capability{core.All},
				},
				RunAsNonRoot:           new(true),
				ReadOnlyRootFileSystem: new(true),
				SELinuxOptions:         &SELinuxOptions{User: "u", Role: "r", Type: "t", Level: "l"},
				SeccompProfile: &SeccompProfile{
					Type:             SeccompProfileTypeLocalhost,
					LocalhostProfile: new("profiles/audit.json"),
				},
				AppArmorProfile: &AppArmorProfile{Type: AppArmorProfileTypeRuntimeDefault},
			},
			SupportMode:         true,
			ExportMode:          true,
			Stopped:             false,
			PauseReconciliation: true,
			UpgradeConfig:       UpgradeConfig{AllowNamespaceSwitch: true, ForceUpgrade: true},
			AdditionalIngressAnnotations: IngressAnnotations{
				"nginx.ingress.kubernetes.io/foo": "bar",
			},
			AdditionalMounts: []DataMount{
				{SourceType: DataSourceConfigMap, Name: "my-cm", Volume: "data", Subfolder: "sub"},
			},
		},
		Status: DoguStatus{
			Status:           DoguStatusInstalled,
			RequeueTime:      5 * time.Second,
			RequeuePhase:     "phase",
			StartedAt:        metav1.NewTime(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)),
			Health:           AvailableHealthStatus,
			InstalledVersion: "1.2.3-4",
			Stopped:          false,
			ExportMode:       true,
			DataVolumeSize:   quantityPtr("5Gi"),
			Conditions: []metav1.Condition{
				{
					Type:               ConditionHealthy,
					Status:             metav1.ConditionTrue,
					Reason:             "Ok",
					Message:            "fine",
					LastTransitionTime: metav1.NewTime(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)),
				},
			},
		},
	}
}

func newCondition(ctype string, status metav1.ConditionStatus, reason, message string, lastTransitionTime metav1.Time, observedGeneration int64) metav1.Condition {
	return metav1.Condition{
		Type:               ctype,
		Status:             status,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: lastTransitionTime,
		ObservedGeneration: observedGeneration,
	}
}

func newTestConditions() []metav1.Condition {
	testTime := metav1.NewTime(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC))
	observingGeneration := int64(1)

	return []metav1.Condition{
		newCondition(v3beta1.ConditionReady, metav1.ConditionFalse, "Installing", "...", testTime, observingGeneration),
		newCondition(v3beta1.ConditionHealthy, metav1.ConditionFalse, "WorkloadsNotReady", "...", testTime, observingGeneration),
		newCondition(v3beta1.ConditionSchemaValidationSkipped, metav1.ConditionFalse, "SchemaValidationEnabled", "...", testTime, observingGeneration),
		newCondition(v3beta1.ConditionValid, metav1.ConditionFalse, "SchemaInvalid", "...", testTime, observingGeneration),
		newCondition(v3beta1.ConditionStopped, metav1.ConditionTrue, "Stopped", "...", testTime, observingGeneration),
		newCondition(v3beta1.ConditionExportModeActive, metav1.ConditionFalse, "NotActive", "...", testTime, observingGeneration),
		newCondition(v3beta1.ConditionPauseReconciliation, metav1.ConditionTrue, "ReconiliationPaused", "...", testTime, observingGeneration),
		newCondition(v3beta1.ConditionChartAvailable, metav1.ConditionFalse, "ChartNotFound", "...", testTime, observingGeneration),
		newCondition(v3beta1.ConditionUpdatePending, metav1.ConditionTrue, "Stopped", "...", testTime, observingGeneration),
	}
}

func newV3beta1TestDogu() *v3beta1.Dogu {
	return &v3beta1.Dogu{
		ObjectMeta: metav1.ObjectMeta{Name: "dogu", Namespace: "ecosystem"},
		Spec: v3beta1.DoguSpec{
			Name:                 "postgresql",
			DoguNamespace:        "official",
			Version:              "1.2.3-4",
			DoguApiVersion:       v3beta1.DoguApiVersionV3,
			Values:               apiextensionsv1.JSON{Raw: []byte(`{"replicas":3}`)},
			MappedValues:         map[string]string{"log-level": "info"},
			SkipSchemaValidation: true,
		},
		Status: v3beta1.DoguStatus{
			AppVersion:       "6.7",
			StartedAt:        metav1.NewTime(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)),
			Conditions:       newTestConditions(),
			InstalledVersion: "1.2.3-4",
			Stopped:          true,
			ExportMode:       false,
			VolumeStatus: []v3beta1.VolumeStatus{
				{
					ClaimName:   "nexus-data",
					Size:        resource.MustParse("10Gi"),
					DesiredSize: resource.MustParse("10Gi"),
					State:       v3beta1.VolumeStatusStateReady,
					MountedBy: []corev1.TypedLocalObjectReference{
						{
							Name:     "nexus",
							Kind:     "StatefulSet",
							APIGroup: new("apps"),
						},
					},
				},
			},
		},
	}
}

func TestDogu_ConvertTo_ConvertFrom_V2NativeRoundTripIsLossless(t *testing.T) {
	src := newFullV2TestDogu()

	var hub v3beta1.Dogu
	require.NoError(t, src.ConvertTo(&hub))

	// a v2-native object (no transport annotations) must default to "v2" and carry no values
	assert.Equal(t, v3beta1.DoguApiVersionV2, hub.Spec.DoguApiVersion)
	assert.Empty(t, hub.Spec.Values.Raw)
	assert.NotContains(t, hub.Annotations, doguApiVersionAnnotationKey)
	assert.NotContains(t, hub.Annotations, valuesAnnotationKey)
	assert.NotContains(t, hub.Annotations, statusAppVersionAnnotationKey)
	assert.NotContains(t, hub.Annotations, mappedValuesAnnotationKey)
	assert.NotContains(t, hub.Annotations, skipSchemaValidationAnnotationKey)
	assert.Equal(t, "unrelated annotation", hub.Annotations["custom.example.com/note"])

	var result Dogu
	require.NoError(t, result.ConvertFrom(&hub))

	assert.Equal(t, src.Spec, result.Spec)
	assert.Equal(t, src.Status, result.Status)
	assert.Equal(t, src.Annotations, result.Annotations)
}

func TestDogu_ConvertFrom_ConvertTo_V3RoundTripIsLossless(t *testing.T) {
	hub := newV3beta1TestDogu()

	var v2View Dogu
	require.NoError(t, v2View.ConvertFrom(hub))

	require.Contains(t, v2View.Annotations, doguApiVersionAnnotationKey)
	require.Contains(t, v2View.Annotations, valuesAnnotationKey)
	require.Contains(t, v2View.Annotations, statusAppVersionAnnotationKey)
	require.Contains(t, v2View.Annotations, mappedValuesAnnotationKey)
	require.Contains(t, v2View.Annotations, skipSchemaValidationAnnotationKey)
	require.Contains(t, v2View.Annotations, statusVolumesAnnotationKey)
	assert.Equal(t, string(v3beta1.DoguApiVersionV3), v2View.Annotations[doguApiVersionAnnotationKey])
	assert.Equal(t, `{"replicas":3}`, v2View.Annotations[valuesAnnotationKey])
	assert.Equal(t, "6.7", v2View.Annotations[statusAppVersionAnnotationKey])
	assert.Equal(t, `{"log-level":"info"}`, v2View.Annotations[mappedValuesAnnotationKey])
	assert.Equal(t, "true", v2View.Annotations[skipSchemaValidationAnnotationKey])
	assert.Equal(t, "official/postgresql", v2View.Spec.Name)
	marshal, err := json.Marshal(hub.Status.VolumeStatus)
	require.NoError(t, err)
	assert.Equal(t, string(marshal), v2View.Annotations[statusVolumesAnnotationKey])

	var restoredHub v3beta1.Dogu
	require.NoError(t, v2View.ConvertTo(&restoredHub))

	assert.Equal(t, v3beta1.DoguApiVersionV3, restoredHub.Spec.DoguApiVersion)
	assert.Equal(t, `{"replicas":3}`, string(restoredHub.Spec.Values.Raw))
	assert.Equal(t, "6.7", restoredHub.Status.AppVersion)
	assert.Equal(t, "official", restoredHub.Spec.DoguNamespace)
	assert.Equal(t, "postgresql", restoredHub.Spec.Name)
	assert.Equal(t, true, restoredHub.Spec.SkipSchemaValidation)
	assert.Equal(t, hub.Status.VolumeStatus, restoredHub.Status.VolumeStatus)
	assert.Equal(t, hub.Status.Conditions, restoredHub.Status.Conditions)
	assert.NotContains(t, restoredHub.Annotations, doguApiVersionAnnotationKey)
	assert.NotContains(t, restoredHub.Annotations, valuesAnnotationKey)
	assert.NotContains(t, restoredHub.Annotations, statusAppVersionAnnotationKey)
	assert.NotContains(t, restoredHub.Annotations, mappedValuesAnnotationKey)
	assert.NotContains(t, restoredHub.Annotations, skipSchemaValidationAnnotationKey)
	assert.NotContains(t, restoredHub.Annotations, statusVolumesAnnotationKey)
}

func TestDogu_ConvertTo_DoesNotMutateSourceAnnotationsMap(t *testing.T) {
	src := &Dogu{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"custom.example.com/note":   "keep me",
				doguApiVersionAnnotationKey: "v3",
				valuesAnnotationKey:         `{"replicas":3}`,
			},
		},
		Spec: DoguSpec{
			Name: "official/postgresql",
		},
	}

	var hub v3beta1.Dogu
	require.NoError(t, src.ConvertTo(&hub))

	// the destination must have the transport keys removed...
	assert.NotContains(t, hub.Annotations, doguApiVersionAnnotationKey)
	assert.NotContains(t, hub.Annotations, valuesAnnotationKey)
	assert.Equal(t, "keep me", hub.Annotations["custom.example.com/note"])

	// ...but the source object must be untouched (no aliasing of the annotations map)
	assert.Equal(t, "v3", src.Annotations[doguApiVersionAnnotationKey])
	assert.Equal(t, `{"replicas":3}`, src.Annotations[valuesAnnotationKey])
	assert.Equal(t, "keep me", src.Annotations["custom.example.com/note"])
}

func TestDogu_ConvertFrom_DoesNotAliasHubAnnotationsMap(t *testing.T) {
	src := &v3beta1.Dogu{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{"custom.example.com/note": "keep me"},
		},
		Spec: v3beta1.DoguSpec{
			DoguApiVersion: v3beta1.DoguApiVersionV3,
			Values:         apiextensionsv1.JSON{Raw: []byte("{}")},
		},
	}

	var dst Dogu
	require.NoError(t, dst.ConvertFrom(src))

	assert.Contains(t, dst.Annotations, doguApiVersionAnnotationKey)
	assert.Contains(t, dst.Annotations, valuesAnnotationKey)

	// the hub's own annotations must not have been mutated by writing the transport keys onto dst
	assert.NotContains(t, src.Annotations, doguApiVersionAnnotationKey)
	assert.NotContains(t, src.Annotations, valuesAnnotationKey)
	assert.Equal(t, "keep me", src.Annotations["custom.example.com/note"])
}

func TestDogu_ConvertTo_ShouldReturnErrorOnInvalidName(t *testing.T) {
	src := &Dogu{
		Spec: DoguSpec{
			Name: "invalid/name/name",
		},
	}

	dst := &v3beta1.Dogu{}
	err := src.ConvertTo(dst)

	require.ErrorContains(t, err, "v2 dogu name \"invalid/name/name\" should only contain one namespace delimiter \"/\"")
}

func TestDogu_ConvertTo_ShouldReturnErrorOnInvalidDoguApiVersionAnnotation(t *testing.T) {
	src := &Dogu{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{doguApiVersionAnnotationKey: "invalid"},
		},
		Spec: DoguSpec{
			Name: "official/postgresql",
		},
	}

	dst := &v3beta1.Dogu{}
	err := src.ConvertTo(dst)

	require.ErrorContains(t, err, "doguApiVersion annotation \"invalid\" is not v2 or v3")
}

func TestDogu_ConvertTo_ShouldReturnErrorOnInvalidValuesAnnotation(t *testing.T) {
	src := &Dogu{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{valuesAnnotationKey: "no: json: string : here"},
		},
		Spec: DoguSpec{
			Name: "official/postgresql",
		},
	}

	dst := &v3beta1.Dogu{}
	err := src.ConvertTo(dst)

	require.ErrorContains(t, err, "values annotation \"no: json: string : here\" is not a valid JSON object")
}

func TestDogu_ConvertTo_ShouldReturnErrorOnInvalidMappedValuesAnnotation(t *testing.T) {
	src := &Dogu{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{mappedValuesAnnotationKey: "no: string: string : here"},
		},
		Spec: DoguSpec{
			Name: "official/postgresql",
		},
	}

	dst := &v3beta1.Dogu{}
	err := src.ConvertTo(dst)

	require.ErrorContains(t, err, "failed to unmarshal mapped values")
}
