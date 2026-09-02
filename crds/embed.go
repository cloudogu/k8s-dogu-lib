package crds

import "embed"

// This embed provides the crd for other applications. They can import this package and use the yaml file
// for the CRD in e.g. integration tests. Otherwise, this file would not be present in the golang vendor directory.
// The file gets refreshed by copying from controller-gen by the "crd-helm-generate/crd-copy-for-go-embedding" make target.
//
//go:embed k8s.cloudogu.com_dogus.yaml
var _ embed.FS
