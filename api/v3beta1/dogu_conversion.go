package v3beta1

import "sigs.k8s.io/controller-runtime/pkg/conversion"

// Compile-Time Interface Assertion
var _ conversion.Hub = &Dogu{}

// Hub marks Dogu as the conversion hub for the k8s.cloudogu.com Dogu kind. Every other served
// API version of Dogu (e.g. v2) implements conversion.Convertible and converts to/from this type.
func (*Dogu) Hub() {}
