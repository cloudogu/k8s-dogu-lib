// Client-gen DOESN'T generate this. YOU CAN EDIT THIS FILE

package v2

import (
	"context"

	apiv2 "github.com/cloudogu/k8s-dogu-lib/v3/api/v2"
	"github.com/cloudogu/retry-lib/retry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DoguExpansion extends the DoguInterface
type DoguExpansion interface {
	UpdateSpecWithRetry(ctx context.Context, dogu *apiv2.Dogu, modifySpecFn func(spec apiv2.DoguSpec) apiv2.DoguSpec, opts metav1.UpdateOptions) (result *apiv2.Dogu, err error)
	UpdateStatusWithRetry(ctx context.Context, dogu *apiv2.Dogu, modifyStatusFn func(apiv2.DoguStatus) apiv2.DoguStatus, opts metav1.UpdateOptions) (result *apiv2.Dogu, err error)
}

// UpdateSpecWithRetry updates the spec of the resource, retrying if a conflict error arises.
func (d *dogus) UpdateSpecWithRetry(ctx context.Context, dogu *apiv2.Dogu, modifySpecFn func(spec apiv2.DoguSpec) apiv2.DoguSpec, opts metav1.UpdateOptions) (result *apiv2.Dogu, err error) {
	firstTry := true

	var currentObj *apiv2.Dogu
	err = retry.OnConflict(func() error {
		if firstTry {
			firstTry = false
			currentObj = dogu.DeepCopy()
		} else {
			currentObj, err = d.Get(ctx, dogu.Name, metav1.GetOptions{})
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

func (d *dogus) UpdateStatusWithRetry(ctx context.Context, dogu *apiv2.Dogu, modifyStatusFn func(apiv2.DoguStatus) apiv2.DoguStatus, opts metav1.UpdateOptions) (result *apiv2.Dogu, err error) {
	firstTry := true

	var currentObj *apiv2.Dogu
	err = retry.OnConflict(func() error {
		if firstTry {
			firstTry = false
			currentObj = dogu.DeepCopy()
		} else {
			currentObj, err = d.Get(ctx, dogu.Name, metav1.GetOptions{})
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
