// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	http "net/http"

	omev1beta1 "github.com/sgl-project/ome/pkg/apis/ome/v1beta1"
	scheme "github.com/sgl-project/ome/pkg/client/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type OmeV1beta1Interface interface {
	RESTClient() rest.Interface
	BaseModelsGetter
	BenchmarkJobsGetter
	ClusterBaseModelsGetter
	ClusterServingRuntimesGetter
	FineTunedWeightsGetter
	InferenceServicesGetter
	ServingRuntimesGetter
}

// OmeV1beta1Client is used to interact with features provided by the ome.io group.
type OmeV1beta1Client struct {
	restClient rest.Interface
}

func (c *OmeV1beta1Client) BaseModels(namespace string) BaseModelInterface {
	return newBaseModels(c, namespace)
}

func (c *OmeV1beta1Client) BenchmarkJobs(namespace string) BenchmarkJobInterface {
	return newBenchmarkJobs(c, namespace)
}

func (c *OmeV1beta1Client) ClusterBaseModels() ClusterBaseModelInterface {
	return newClusterBaseModels(c)
}

func (c *OmeV1beta1Client) ClusterServingRuntimes() ClusterServingRuntimeInterface {
	return newClusterServingRuntimes(c)
}

func (c *OmeV1beta1Client) FineTunedWeights() FineTunedWeightInterface {
	return newFineTunedWeights(c)
}

func (c *OmeV1beta1Client) InferenceServices(namespace string) InferenceServiceInterface {
	return newInferenceServices(c, namespace)
}

func (c *OmeV1beta1Client) ServingRuntimes(namespace string) ServingRuntimeInterface {
	return newServingRuntimes(c, namespace)
}

// NewForConfig creates a new OmeV1beta1Client for the given config.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*OmeV1beta1Client, error) {
	config := *c
	setConfigDefaults(&config)
	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(&config, httpClient)
}

// NewForConfigAndClient creates a new OmeV1beta1Client for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
func NewForConfigAndClient(c *rest.Config, h *http.Client) (*OmeV1beta1Client, error) {
	config := *c
	setConfigDefaults(&config)
	client, err := rest.RESTClientForConfigAndClient(&config, h)
	if err != nil {
		return nil, err
	}
	return &OmeV1beta1Client{client}, nil
}

// NewForConfigOrDie creates a new OmeV1beta1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *OmeV1beta1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new OmeV1beta1Client for the given RESTClient.
func New(c rest.Interface) *OmeV1beta1Client {
	return &OmeV1beta1Client{c}
}

func setConfigDefaults(config *rest.Config) {
	gv := omev1beta1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = rest.CodecFactoryForGeneratedClient(scheme.Scheme, scheme.Codecs).WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *OmeV1beta1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
