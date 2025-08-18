package plugins

import (
	"context"

	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	extensionsplug "github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugin"
	"github.com/kgateway-dev/kgateway/v2/pkg/settings"
)

type AgwCollections struct {
	Client kube.Client

	// common k8s collections
	Services krt.Collection[*corev1.Service]

	// common kgateway collections
	GatewayExtensions krt.Collection[*v1alpha1.GatewayExtension]

	ControllerName string
}

func (c *AgwCollections) HasSynced() bool {
	// we check nil as well because some of the inner
	// collections aren't initialized until we call InitPlugins
	return c.GatewayExtensions != nil && c.GatewayExtensions.HasSynced() &&
		c.Services != nil && c.Services.HasSynced()
}

// NewAgwCollections initializes the core krt collections.
// Collections that rely on plugins aren't initialized here,
// and InitPlugins must be called.
func NewAgwCollections(
	client kube.Client,
	controllerName string,
	services krt.Collection[*corev1.Service],
	gwExts krt.Collection[*v1alpha1.GatewayExtension],
) (*AgwCollections, error) {

	return &AgwCollections{
		Client:            client,
		Services:          services,
		GatewayExtensions: gwExts,
		ControllerName:    controllerName,
	}, nil
}

// InitPlugins set up collections that rely on plugins.
// This can't be part of NewAgwCollections because the setup
// of plugins themselves rely on a reference to AgwCollections.
func (c *AgwCollections) InitPlugins(
	ctx context.Context,
	mergedPlugins extensionsplug.Plugin,
	globalSettings settings.Settings,
) {
	// TODO: move from new to init
}
