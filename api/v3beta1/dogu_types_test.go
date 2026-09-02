package v3beta1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testDogu = &Dogu{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "k8s.cloudogu.com/v3beta1",
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

func TestDogu_GetObjectMeta(t *testing.T) {
	actual := testDogu.GetObjectMeta()

	expectedObjKey := &metav1.ObjectMeta{
		Namespace: "ecosystem",
		Name:      "dogu",
	}
	assert.Equal(t, expectedObjKey, actual)
}

func TestDogu_GetSimpleNameVersion(t *testing.T) {
	t.Run("should get simple name version", func(t *testing.T) {
		sut := &Dogu{ObjectMeta: metav1.ObjectMeta{Name: "dogu"}, Spec: DoguSpec{Name: "dogu", Version: "1.2.3"}}
		version, err := sut.GetSimpleNameVersion()
		assert.NoError(t, err)
		assert.Equal(t, "dogu", version.Name.String())
		assert.Equal(t, "1.2.3", version.Version.String())
	})

	t.Run("should fail to parse version", func(t *testing.T) {
		sut := &Dogu{ObjectMeta: metav1.ObjectMeta{Name: "dogu"}, Spec: DoguSpec{Name: "dogu", Version: "1.2xxx"}}
		_, err := sut.GetSimpleNameVersion()
		assert.Error(t, err)
	})

}

func TestDogu_GetConditions(t *testing.T) {
	dogu := Dogu{
		Status: DoguStatus{
			Conditions: []metav1.Condition{
				{Type: "aType"},
			},
		},
	}

	conditions := dogu.GetConditions()

	assert.Equal(t, "aType", conditions[0].Type)
}

func TestDogu_SetConditions(t *testing.T) {
	dogu := Dogu{}

	dogu.SetConditions([]metav1.Condition{
		{Type: "aType"},
	})

	assert.Equal(t, "aType", dogu.Status.Conditions[0].Type)

}
