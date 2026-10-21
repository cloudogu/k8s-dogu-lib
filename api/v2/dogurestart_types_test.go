package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDoguRestart_GetConditions(t *testing.T) {
	t.Run("should return conditions from status", func(t *testing.T) {
		// given
		expected := []metav1.Condition{{Type: DoguRestartConditionReady, Status: metav1.ConditionTrue, Reason: DoguRestartReasonSucceeded}}
		sut := &DoguRestart{Status: DoguRestartStatus{Conditions: expected}}

		// when
		actual := sut.GetConditions()

		// then
		assert.Equal(t, expected, actual)
	})
}

func TestDoguRestart_SetConditions(t *testing.T) {
	t.Run("should set conditions in status without touching the phase", func(t *testing.T) {
		// given
		sut := &DoguRestart{}
		sut.Status.Phase = RestartStatusPhaseStarting
		conditions := []metav1.Condition{
			{Type: DoguRestartConditionRestarted, Status: metav1.ConditionTrue, Reason: DoguRestartReasonSucceeded},
			{Type: DoguRestartConditionReady, Status: metav1.ConditionFalse, Reason: DoguRestartReasonRolloutInProgress},
		}

		// when
		sut.SetConditions(conditions)

		// then
		assert.Equal(t, conditions, sut.Status.Conditions)
		assert.Equal(t, RestartStatusPhaseStarting, sut.Status.Phase)
	})

	t.Run("should work with meta.SetStatusCondition", func(t *testing.T) {
		// given
		sut := &DoguRestart{}
		conditions := sut.GetConditions()

		// when
		meta.SetStatusCondition(&conditions, metav1.Condition{Type: DoguRestartConditionReady, Status: metav1.ConditionFalse, Reason: DoguRestartReasonDoguNotFound})
		sut.SetConditions(conditions)

		// then
		cond := meta.FindStatusCondition(sut.GetConditions(), DoguRestartConditionReady)
		require.NotNil(t, cond)
		assert.Equal(t, DoguRestartReasonDoguNotFound, cond.Reason)
		assert.Equal(t, cond.Status, metav1.ConditionFalse)
	})
}
