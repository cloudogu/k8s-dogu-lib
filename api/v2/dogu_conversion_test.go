package v2

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/cesapp-lib/core"
	v3beta1 "github.com/cloudogu/k8s-dogu-lib/v2/api/v3beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	assert.Equal(t, "unrelated annotation", hub.Annotations["custom.example.com/note"])

	var result Dogu
	require.NoError(t, result.ConvertFrom(&hub))

	assert.Equal(t, src.Spec, result.Spec)
	assert.Equal(t, src.Status, result.Status)
	assert.Equal(t, src.Annotations, result.Annotations)
}

func TestDogu_ConvertFrom_ConvertTo_V3RoundTripIsLossless(t *testing.T) {
	hub := &v3beta1.Dogu{
		ObjectMeta: metav1.ObjectMeta{Name: "dogu", Namespace: "ecosystem"},
		Spec: v3beta1.DoguSpec{
			Name:           "postgresql",
			DoguNamespace:  "official",
			Version:        "1.2.3-4",
			DoguApiVersion: v3beta1.DoguApiVersionV3,
			Values:         runtime.RawExtension{Raw: []byte(`{"replicas":3}`)},
		},
		Status: v3beta1.DoguStatus{
			AppVersion: "1.2.3-4",
		},
	}

	var v2View Dogu
	require.NoError(t, v2View.ConvertFrom(hub))

	require.Contains(t, v2View.Annotations, doguApiVersionAnnotationKey)
	require.Contains(t, v2View.Annotations, valuesAnnotationKey)
	require.Contains(t, v2View.Annotations, statusAppVersionAnnotationKey)
	assert.Equal(t, string(v3beta1.DoguApiVersionV3), v2View.Annotations[doguApiVersionAnnotationKey])
	assert.Equal(t, `{"replicas":3}`, v2View.Annotations[valuesAnnotationKey])
	assert.Equal(t, "1.2.3-4", v2View.Annotations[statusAppVersionAnnotationKey])
	assert.Equal(t, "official/postgresql", v2View.Spec.Name)

	var restoredHub v3beta1.Dogu
	require.NoError(t, v2View.ConvertTo(&restoredHub))

	assert.Equal(t, v3beta1.DoguApiVersionV3, restoredHub.Spec.DoguApiVersion)
	assert.Equal(t, `{"replicas":3}`, string(restoredHub.Spec.Values.Raw))
	assert.Equal(t, "1.2.3-4", restoredHub.Status.AppVersion)
	assert.Equal(t, "official", restoredHub.Spec.DoguNamespace)
	assert.Equal(t, "postgresql", restoredHub.Spec.Name)
	assert.NotContains(t, restoredHub.Annotations, doguApiVersionAnnotationKey)
	assert.NotContains(t, restoredHub.Annotations, valuesAnnotationKey)
	assert.NotContains(t, restoredHub.Annotations, statusAppVersionAnnotationKey)
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
			Values:         runtime.RawExtension{Raw: []byte("{}")},
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
