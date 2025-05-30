package v2

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/cesapp-lib/core"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	eventV1 "k8s.io/api/events/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var testDogu = &Dogu{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "k8s.cloudogu.com/v2",
		Kind:       "Dogu",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "dogu",
		Namespace: "ecosystem",
	},
	Spec: DoguSpec{
		Name:          "namespace/dogu",
		Version:       "1.2.3-4",
		UpgradeConfig: UpgradeConfig{},
	},
	Status: DoguStatus{Status: ""},
}
var testCtx = context.Background()

func TestDoguStatus_GetRequeueTime(t *testing.T) {
	tests := []struct {
		requeueCount        time.Duration
		expectedRequeueTime time.Duration
	}{
		{requeueCount: time.Second, expectedRequeueTime: time.Second * 2},
		{requeueCount: time.Second * 17, expectedRequeueTime: time.Second * 34},
		{requeueCount: time.Minute, expectedRequeueTime: time.Minute * 2},
		{requeueCount: time.Minute * 7, expectedRequeueTime: time.Minute * 14},
		{requeueCount: time.Minute * 45, expectedRequeueTime: time.Hour*1 + time.Minute*30},
		{requeueCount: time.Hour * 2, expectedRequeueTime: time.Hour * 4},
		{requeueCount: time.Hour * 3, expectedRequeueTime: time.Hour * 6},
		{requeueCount: time.Hour * 5, expectedRequeueTime: time.Hour * 6},
		{requeueCount: time.Hour * 100, expectedRequeueTime: time.Hour * 6},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("calculate next requeue time for current time %s", tt.requeueCount), func(t *testing.T) {
			// given
			ds := &DoguStatus{
				RequeueTime: tt.requeueCount,
			}

			// when
			actualRequeueTime := ds.NextRequeue()

			// then
			assert.Equal(t, tt.expectedRequeueTime, actualRequeueTime)
		})
	}
}

func TestDoguStatus_ResetRequeueTime(t *testing.T) {
	t.Run("reset requeue time", func(t *testing.T) {
		// given
		ds := &DoguStatus{
			RequeueTime: time.Hour * 3,
		}

		// when
		ds.ResetRequeueTime()

		// then
		assert.Equal(t, RequeueTimeInitialRequeueTime, ds.RequeueTime)
	})
}

func TestDogu_GetSecretObjectKey(t *testing.T) {
	// given
	ds := &Dogu{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myspecialdogu",
			Namespace: "testnamespace",
		},
	}

	// when
	key := ds.GetSecretObjectKey()

	// then
	assert.Equal(t, "myspecialdogu-secrets", key.Name)
	assert.Equal(t, "testnamespace", key.Namespace)
}

func TestDogu_GetObjectKey(t *testing.T) {
	actual := testDogu.GetObjectKey()

	expectedObjKey := client.ObjectKey{
		Namespace: "ecosystem",
		Name:      "dogu",
	}
	assert.Equal(t, expectedObjKey, actual)
}

func TestDogu_GetObjectMeta(t *testing.T) {
	actual := testDogu.GetObjectMeta()

	expectedObjKey := &metav1.ObjectMeta{
		Namespace: "ecosystem",
		Name:      "dogu",
	}
	assert.Equal(t, expectedObjKey, actual)
}

func TestDogu_GetDataVolumeName(t *testing.T) {
	actual := testDogu.GetDataVolumeName()

	assert.Equal(t, "dogu-data", actual)
}

func TestDogu_GetPrivateVolumeName(t *testing.T) {
	actual := testDogu.GetPrivateKeySecretName()

	assert.Equal(t, "dogu-private", actual)
}

func TestDogu_GetDevelopmentDoguMapKey(t *testing.T) {
	actual := testDogu.GetDevelopmentDoguMapKey()

	expectedKey := client.ObjectKey{
		Namespace: "ecosystem",
		Name:      "dogu-descriptor",
	}
	assert.Equal(t, expectedKey, actual)
}

func TestDogu_GetPrivateKeyObjectKey(t *testing.T) {
	actual := testDogu.GetPrivateKeyObjectKey()

	expectedKey := client.ObjectKey{
		Namespace: "ecosystem",
		Name:      "dogu-private",
	}
	assert.Equal(t, expectedKey, actual)
}

func TestCesMatchingLabels_Add(t *testing.T) {
	t.Run("should add to empty object", func(t *testing.T) {
		input := CesMatchingLabels{"key": "value"}
		// when
		actual := CesMatchingLabels{}.Add(input)

		// then
		require.NotEmpty(t, actual)
		expected := CesMatchingLabels{"key": "value"}
		assert.Equal(t, expected, actual)
	})
	t.Run("should add to filed object", func(t *testing.T) {
		input := CesMatchingLabels{"key2": "value2"}
		// when
		actual := CesMatchingLabels{"key": "value"}.Add(input)

		// then
		require.NotEmpty(t, actual)
		expected := CesMatchingLabels{"key": "value", "key2": "value2"}
		assert.Equal(t, expected, actual)
	})
}

func TestDogu_Labels(t *testing.T) {
	sut := Dogu{
		ObjectMeta: metav1.ObjectMeta{Namespace: "test", Name: "ldap"},
		Spec: DoguSpec{
			Name:    "official/ldap",
			Version: "1.2.3-4",
		},
	}

	t.Run("should return pod labels", func(t *testing.T) {
		actual := sut.GetPodLabels()

		expected := CesMatchingLabels{"dogu.name": "ldap", "dogu.version": "1.2.3-4"}
		assert.Equal(t, expected, actual)
	})

	t.Run("should return dogu name label", func(t *testing.T) {
		// when
		actual := sut.GetDoguNameLabel()

		// then
		expected := CesMatchingLabels{"dogu.name": "ldap"}
		assert.Equal(t, expected, actual)
	})
}

func TestDogu_GetPod(t *testing.T) {
	readyPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "ldap-x2y3z45", Labels: testDogu.GetPodLabels()},
		Status:     corev1.PodStatus{Conditions: []corev1.PodCondition{{Type: corev1.ContainersReady, Status: corev1.ConditionTrue}}},
	}
	cli := fake.NewClientBuilder().WithScheme(getDoguTypesTestScheme()).WithObjects(readyPod).Build()

	// when
	actual, err := testDogu.GetPod(testCtx, cli)

	// then
	require.NoError(t, err)
	exptectedPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "ldap-x2y3z45", Labels: testDogu.GetPodLabels()},
		Status:     corev1.PodStatus{Conditions: []corev1.PodCondition{{Type: corev1.ContainersReady, Status: corev1.ConditionTrue}}},
	}
	// ignore ResourceVersion which is introduced during getting pods from the K8s API
	actual.ResourceVersion = ""
	assert.Equal(t, exptectedPod, actual)
}

func TestDevelopmentDoguMap_DeleteFromCluster(t *testing.T) {
	t.Run("should delete a DevelopmentDogu cm", func(t *testing.T) {
		inputCm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "ldap-dev-dev-map"},
			Data:       map[string]string{"key": "le data"},
		}
		mockClient := newMockK8sClient(t)
		mockClient.EXPECT().Delete(testCtx, inputCm).Return(nil)
		sut := &DevelopmentDoguMap{
			ObjectMeta: metav1.ObjectMeta{Name: "ldap-dev-dev-map"},
			Data:       map[string]string{"key": "le data"},
		}

		// when
		err := sut.DeleteFromCluster(testCtx, mockClient)

		// then
		require.NoError(t, err)
	})
	t.Run("should return an error", func(t *testing.T) {
		inputCm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "ldap-dev-dev-map"},
			Data:       map[string]string{"key": "le data"},
		}
		mockClient := newMockK8sClient(t)
		mockClient.EXPECT().Delete(testCtx, inputCm).Return(assert.AnError)
		sut := &DevelopmentDoguMap{
			ObjectMeta: metav1.ObjectMeta{Name: "ldap-dev-dev-map"},
			Data:       map[string]string{"key": "le data"},
		}

		// when
		err := sut.DeleteFromCluster(testCtx, mockClient)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func getDoguTypesTestScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()

	scheme.AddKnownTypeWithName(schema.GroupVersionKind{
		Group:   "k8s.cloudogu.com",
		Version: "v2",
		Kind:    "dogu",
	}, &Dogu{})
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "Deployment",
	}, &appsv1.Deployment{})
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Secret",
	}, &corev1.Secret{})
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Service",
	}, &corev1.Service{})
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "PersistentVolumeClaim",
	}, &corev1.PersistentVolumeClaim{})
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "ConfigMap",
	}, &corev1.ConfigMap{})
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Event",
	}, &eventV1.Event{})
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Pod",
	}, &corev1.Pod{})
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "PodList",
	}, &corev1.PodList{})

	return scheme
}

func TestDogu_GetPrivateKeySecret(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		expected := &corev1.Secret{
			TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
			ObjectMeta: metav1.ObjectMeta{Name: "dogu-private", Namespace: "ecosystem"},
		}
		fakeClient := fake.NewClientBuilder().WithScheme(getDoguTypesTestScheme()).WithObjects(expected).Build()

		// when
		secret, err := testDogu.GetPrivateKeySecret(context.TODO(), fakeClient)

		// then
		require.NoError(t, err)
		assert.Equal(t, expected, secret)
	})

	t.Run("fail to get private key secret", func(t *testing.T) {
		// given
		fakeClient := fake.NewClientBuilder().WithScheme(getDoguTypesTestScheme()).Build()

		// when
		_, err := testDogu.GetPrivateKeySecret(context.TODO(), fakeClient)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to get private key secret for dogu")
	})
}

func TestDogu_ChangeRequeuePhaseWithRetry(t *testing.T) {
	t.Run("success on conflict", func(t *testing.T) {
		// given
		resourceVersion := "1"
		sut := &Dogu{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "postgresql",
				Namespace:       "ecosystem",
				ResourceVersion: resourceVersion,
			},
		}

		newDogu := &Dogu{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "postgresql",
				Namespace:       "ecosystem",
				ResourceVersion: "2",
			},
			Status: DoguStatus{
				RequeuePhase: "old",
			},
		}

		requeuePhase := "phase"
		fakeClient := fake.NewClientBuilder().WithScheme(getDoguTypesTestScheme()).WithStatusSubresource(&Dogu{}).WithObjects(newDogu).Build()

		// when
		err := sut.ChangeRequeuePhaseWithRetry(testCtx, fakeClient, requeuePhase)

		// then
		require.NoError(t, err)
		assert.Equal(t, requeuePhase, sut.Status.RequeuePhase)
		assert.NotEqual(t, resourceVersion, sut.ResourceVersion)
	})

	t.Run("should return error on get error", func(t *testing.T) {
		// given
		resourceVersion := "1"
		sut := &Dogu{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "postgresql",
				Namespace:       "ecosystem",
				ResourceVersion: resourceVersion,
			},
		}

		requeuePhase := "phase"
		fakeClient := fake.NewClientBuilder().WithScheme(getDoguTypesTestScheme()).Build()

		// when
		err := sut.ChangeRequeuePhaseWithRetry(testCtx, fakeClient, requeuePhase)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "dogus.k8s.cloudogu.com \"postgresql\" not found")
	})
}

func TestDogu_ChangeStateWithRetry(t *testing.T) {
	t.Run("success on conflict", func(t *testing.T) {
		// given
		resourceVersion := "1"
		sut := &Dogu{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "postgresql",
				Namespace:       "ecosystem",
				ResourceVersion: resourceVersion,
			},
		}

		newDogu := &Dogu{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "postgresql",
				Namespace:       "ecosystem",
				ResourceVersion: "2",
			},
			Status: DoguStatus{
				Status: "old",
			},
		}

		status := "status"
		fakeClient := fake.NewClientBuilder().WithScheme(getDoguTypesTestScheme()).WithStatusSubresource(&Dogu{}).WithObjects(newDogu).Build()

		// when
		err := sut.ChangeStateWithRetry(testCtx, fakeClient, status)

		// then
		require.NoError(t, err)
		assert.Equal(t, status, sut.Status.Status)
		assert.NotEqual(t, resourceVersion, sut.ResourceVersion)
	})

	t.Run("should return error on get error", func(t *testing.T) {
		// given
		resourceVersion := "1"
		sut := &Dogu{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "postgresql",
				Namespace:       "ecosystem",
				ResourceVersion: resourceVersion,
			},
		}

		status := "status"
		fakeClient := fake.NewClientBuilder().WithScheme(getDoguTypesTestScheme()).Build()

		// when
		err := sut.ChangeStateWithRetry(testCtx, fakeClient, status)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "dogus.k8s.cloudogu.com \"postgresql\" not found")
	})
}

func TestDogu_NextRequeueWithRetry(t *testing.T) {
	t.Run("success on conflict; requeue time was reset", func(t *testing.T) {
		// given
		sut := Dogu{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "postgresql",
				Namespace:       "ecosystem",
				ResourceVersion: "1",
			},
			Status: DoguStatus{
				RequeueTime: time.Second * 40,
			},
		}

		newDogu := &Dogu{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "postgresql",
				Namespace:       "ecosystem",
				ResourceVersion: "2",
			},
			Status: DoguStatus{
				RequeueTime: 0,
			},
		}
		fakeClient := fake.NewClientBuilder().WithScheme(getDoguTypesTestScheme()).WithObjects(newDogu).WithStatusSubresource(&Dogu{}).Build()

		// when
		retry, err := sut.NextRequeueWithRetry(testCtx, fakeClient)

		// then
		require.NoError(t, err)
		assert.Equal(t, time.Second*10, retry)
	})

	t.Run("should return error on get error", func(t *testing.T) {
		// given
		sut := &Dogu{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "postgresql",
				Namespace: "ecosystem",
			},
		}

		fakeClient := fake.NewClientBuilder().WithScheme(getDoguTypesTestScheme()).Build()

		// when
		_, err := sut.NextRequeueWithRetry(testCtx, fakeClient)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "dogus.k8s.cloudogu.com \"postgresql\" not found")
	})
}

func TestDogu_ValidateSecurity(t *testing.T) {
	type args struct {
		dogu *Dogu
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{"valid empty", args{&Dogu{}}, assert.NoError},
		{"valid add filled", args{&Dogu{Spec: DoguSpec{Security: Security{Capabilities: Capabilities{Add: []core.Capability{core.AuditControl}}}}}}, assert.NoError},
		{"valid drop filled", args{&Dogu{Spec: DoguSpec{Security: Security{Capabilities: Capabilities{Drop: []core.Capability{core.AuditControl}}}}}}, assert.NoError},
		{"all possible values", args{&Dogu{Spec: DoguSpec{Security: Security{Capabilities: Capabilities{Add: core.AllCapabilities, Drop: core.AllCapabilities}}}}}, assert.NoError},
		{"add all keyword", args{&Dogu{Spec: DoguSpec{Security: Security{Capabilities: Capabilities{Add: []core.Capability{core.All}}}}}}, assert.NoError},
		{"drop all keyword", args{&Dogu{Spec: DoguSpec{Security: Security{Capabilities: Capabilities{Drop: []core.Capability{core.All}}}}}}, assert.NoError},

		{"invalid valid add filled", args{&Dogu{Spec: DoguSpec{Security: Security{Capabilities: Capabilities{Add: []core.Capability{"err"}}}}}}, assert.Error},
		{"invalid valid drop filled", args{&Dogu{Spec: DoguSpec{Security: Security{Capabilities: Capabilities{Drop: []core.Capability{"err"}}}}}}, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, tt.args.dogu.ValidateSecurity(), fmt.Sprintf("ValidateSecurity(%v)", tt.args.dogu))
		})
	}
}

func TestDogu_ValidateSecurity_message(t *testing.T) {
	t.Run("should match for drop errors", func(t *testing.T) {
		// given
		dogu := &Dogu{Spec: DoguSpec{Name: "official/dogu", Version: "1.2.3", Security: Security{Capabilities: Capabilities{Drop: []core.Capability{"err"}}}}}

		// when
		actual := dogu.ValidateSecurity()

		// then
		require.Error(t, actual)
		assert.ErrorContains(t, actual, "dogu resource official/dogu:1.2.3 contains at least one invalid security field: err is not a valid capability to be dropped")
	})
	t.Run("should match for add errors", func(t *testing.T) {
		// given
		dogu := &Dogu{Spec: DoguSpec{Name: "official/dogu", Version: "1.2.3", Security: Security{Capabilities: Capabilities{Add: []core.Capability{"err"}}}}}

		// when
		actual := dogu.ValidateSecurity()

		// then
		require.Error(t, actual)
		assert.ErrorContains(t, actual, "dogu resource official/dogu:1.2.3 contains at least one invalid security field: err is not a valid capability to be added")
	})
}

func TestDogu_GetMinDataVolumeSize(t *testing.T) {

	testQuantity := func(minSize *resource.Quantity, size *string, expected int64, errortext *string) {
		var dogu *Dogu
		dogu = &Dogu{
			Spec: DoguSpec{
				Resources: DoguResources{},
			},
		}

		if minSize != nil {
			dogu.Spec.Resources.MinDataVolumeSize = minSize.DeepCopy()
		}
		if size != nil {
			dogu.Spec.Resources.DataVolumeSize = *size
		}

		// when
		actual, err := dogu.GetMinDataVolumeSize()

		// then
		if errortext == nil {
			require.NoError(t, err)
			assert.Equal(t, expected, actual.Value())
		} else {
			require.Error(t, err)
			assert.True(t, actual.IsZero())
			assert.ErrorContains(t, err, *errortext)
		}
	}

	t.Run("min Data volume size should be default", func(t *testing.T) {
		testQuantity(nil, nil, int64(2147483648), nil)
	})
	t.Run("min Data volume size should be 1Gi", func(t *testing.T) {
		q, err := resource.ParseQuantity("1Gi")
		require.NoError(t, err)
		testQuantity(&q, nil, int64(1073741824), nil)
	})
	t.Run("min Data volume size should be set to zero so default is returned", func(t *testing.T) {
		q, err := resource.ParseQuantity("0")
		require.NoError(t, err)
		testQuantity(&q, nil, int64(2147483648), nil)
	})
	t.Run("Data volume size should be returned as fallback for empty min data volume size", func(t *testing.T) {
		minsize, err := resource.ParseQuantity("0")
		require.NoError(t, err)
		size := "3Gi"
		testQuantity(&minsize, &size, int64(3221225472), nil)
	})
	t.Run("parsing data volume size should fail", func(t *testing.T) {
		minsize, err := resource.ParseQuantity("0")
		require.NoError(t, err)
		size := "invalid"
		errorText := "quantities must match the regular expression"
		testQuantity(&minsize, &size, int64(3221225472), &errorText)
	})
}
