package client

import (
	"github.com/cloudogu/k8s-dogu-lib/v2/api/v3beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type EcoSystemV3beta1Interface interface {
	Dogus(namespace string) DoguV3beta1Interface
}

type EcoSystemV3beta1Client struct {
	restClient rest.Interface
}

func NewForConfigV3beta1(c *rest.Config) (*EcoSystemV3beta1Client, error) {
	config := *c
	gv := schema.GroupVersion{Group: v3beta1.GroupVersion.Group, Version: v3beta1.GroupVersion.Version}
	config.ContentConfig.GroupVersion = &gv
	config.APIPath = "/apis"

	s := scheme.Scheme
	err := v3beta1.AddToScheme(s)
	if err != nil {
		return nil, err
	}

	metav1.AddToGroupVersion(s, gv)
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &EcoSystemV3beta1Client{restClient: client}, nil
}

func (c *EcoSystemV3beta1Client) Dogus(namespace string) DoguV3beta1Interface {
	return &doguV3beta1Client{
		client: c.restClient,
		ns:     namespace,
	}
}
