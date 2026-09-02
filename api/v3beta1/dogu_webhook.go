package v3beta1

import ctrl "sigs.k8s.io/controller-runtime"

// SetupWebhookWithManager registers the CRD conversion webhook for the Dogu kind with mgr. Since
// Dogu implements conversion.Hub (see dogu_conversion.go), this mounts the generic conversion
// handler that lets the apiserver convert between v3beta1 and any other served, Convertible
// spoke version (currently v2). k8s-dogu-operator must call this exactly once from its main.go,
// after registering v3beta1.AddToScheme on the manager's scheme:
//
//	if err := (&v3beta1.Dogu{}).SetupWebhookWithManager(mgr); err != nil { ... }
func (d *Dogu) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, d).Complete()
}
