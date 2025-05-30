package v2

import (
	"context"
	"testing"

	cloudoguerrors "github.com/cloudogu/ces-commons-lib/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	eventV1 "k8s.io/api/events/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGetPodForLabels(t *testing.T) {
	var testCtx = context.Background()

	t.Run("should return a pod for given labels", func(t *testing.T) {
		// given
		labels := CesMatchingLabels{DoguLabelName: "ldap", DoguLabelVersion: "1.2.3-4"}
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "ldap-x2y3z45", Labels: labels},
			Status:     corev1.PodStatus{Conditions: []corev1.PodCondition{{Type: corev1.ContainersReady, Status: corev1.ConditionTrue}}},
		}
		cli := fake.NewClientBuilder().WithScheme(getPodFinderTestScheme()).WithObjects(pod).Build()

		// when
		actual, err := GetPodForLabels(testCtx, cli, labels)

		// then
		require.NoError(t, err)

		// ignore ResourceVersion which is introduced during getting pods from the K8s API
		actual.ResourceVersion = ""
		expected := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "ldap-x2y3z45", Labels: labels},
			Status:     corev1.PodStatus{Conditions: []corev1.PodCondition{{Type: corev1.ContainersReady, Status: corev1.ConditionTrue}}},
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should return an when no pod was found", func(t *testing.T) {
		// given
		labels := CesMatchingLabels{DoguLabelName: "ldap", DoguLabelVersion: "1.2.3-4"}
		cli := newMockK8sClient(t)
		cli.On("List", testCtx, mock.Anything, client.MatchingLabels(labels)).Return(assert.AnError)

		// when
		_, err := GetPodForLabels(testCtx, cli, labels)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to get pods")
	})
	t.Run("should return an error when no pod was found", func(t *testing.T) {
		// given
		labels := CesMatchingLabels{DoguLabelName: "ldap", DoguLabelVersion: "1.2.3-4"}
		cli := newMockK8sClient(t)
		cli.On("List", testCtx, mock.Anything, client.MatchingLabels(labels)).Return(nil)

		// when
		_, err := GetPodForLabels(testCtx, cli, labels)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "found no pods for labels")
		assert.True(t, cloudoguerrors.IsNotFoundError(err))
	})
	t.Run("should return for multiple pods for unique labels", func(t *testing.T) {
		// given
		labels := CesMatchingLabels{DoguLabelName: "ldap", DoguLabelVersion: "1.2.3-4"}
		pod1 := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "ldap-1", Labels: labels},
			Status:     corev1.PodStatus{Conditions: []corev1.PodCondition{{Type: corev1.ContainersReady, Status: corev1.ConditionTrue}}},
		}
		pod2 := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "ldap-2", Labels: labels},
			Status:     corev1.PodStatus{Conditions: []corev1.PodCondition{{Type: corev1.ContainersReady, Status: corev1.ConditionTrue}}},
		}
		cli := fake.NewClientBuilder().WithScheme(getPodFinderTestScheme()).WithObjects(pod1, pod2).Build()

		// when
		_, err := GetPodForLabels(testCtx, cli, labels)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "found more than one pod")
	})
}

func getPodFinderTestScheme() *runtime.Scheme {
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
