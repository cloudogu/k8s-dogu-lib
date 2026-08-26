// Client-gen DOESN'T generate this. YOU CAN EDIT THIS FILE

package v2

import (
	"context"

	apiv2 "github.com/cloudogu/k8s-dogu-lib/v2/api/v2"
	"github.com/cloudogu/retry-lib/retry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DoguRestartExpansion extends the DoguRestartInterface
type DoguRestartExpansion interface {
	UpdateSpecWithRetry(ctx context.Context, doguRestart *apiv2.DoguRestart, modifySpecFn func(spec apiv2.DoguRestartSpec) apiv2.DoguRestartSpec, opts metav1.UpdateOptions) (result *apiv2.DoguRestart, err error)
	UpdateStatusWithRetry(ctx context.Context, doguRestart *apiv2.DoguRestart, modifyStatusFn func(apiv2.DoguRestartStatus) apiv2.DoguRestartStatus, opts metav1.UpdateOptions) (result *apiv2.DoguRestart, err error)
}

// UpdateSpecWithRetry updates the spec of the resource, retrying if a conflict error arises.
func (d doguRestarts) UpdateSpecWithRetry(ctx context.Context, doguRestart *apiv2.DoguRestart, modifySpecFn func(spec apiv2.DoguRestartSpec) apiv2.DoguRestartSpec, opts metav1.UpdateOptions) (result *apiv2.DoguRestart, err error) {
	firstTry := true

	var currentObj *apiv2.DoguRestart
	err = retry.OnConflict(func() error {
		if firstTry {
			firstTry = false
			currentObj = doguRestart.DeepCopy()
		} else {
			currentObj, err = d.Get(ctx, doguRestart.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
		}

		currentObj.Spec = modifySpecFn(currentObj.Spec)
		currentObj, err = d.Update(ctx, currentObj, opts)
		return err
	})
	if err != nil {
		return nil, err
	}

	return currentObj, nil
}

// UpdateStatusWithRetry updates the status of the resource, retrying if a conflict error arises.
func (d doguRestarts) UpdateStatusWithRetry(ctx context.Context, doguRestart *apiv2.DoguRestart, modifyStatusFn func(apiv2.DoguRestartStatus) apiv2.DoguRestartStatus, opts metav1.UpdateOptions) (result *apiv2.DoguRestart, err error) {
	firstTry := true

	var currentObj *apiv2.DoguRestart
	err = retry.OnConflict(func() error {
		if firstTry {
			firstTry = false
			currentObj = doguRestart.DeepCopy()
		} else {
			currentObj, err = d.Get(ctx, doguRestart.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
		}

		currentObj.Status = modifyStatusFn(currentObj.Status)
		currentObj, err = d.UpdateStatus(ctx, currentObj, opts)
		return err
	})
	if err != nil {
		return nil, err
	}

	return currentObj, nil
}
