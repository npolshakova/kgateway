package agentgatewaysyncer

import (
	"istio.io/istio/pkg/kube/krt"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gateway "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
)

type GatewayClass struct {
	Name       string
	Controller gateway.GatewayController
}

func (g GatewayClass) ResourceName() string {
	return g.Name
}

func GatewayClassesCollection(
	gatewayClasses krt.Collection[*gateway.GatewayClass],
	krtopts krtutil.KrtOptions,
) krt.Collection[GatewayClass] {
	return krt.NewCollection(gatewayClasses, func(ctx krt.HandlerContext, obj *gateway.GatewayClass) *GatewayClass {
		return &GatewayClass{
			Name:       obj.Name,
			Controller: obj.Spec.ControllerName,
		}
	}, krtopts.ToOptions("GatewayClasses")...)
}

func fetchClass(ctx krt.HandlerContext, gatewayClasses krt.Collection[GatewayClass], gc gatewayv1.ObjectName) *GatewayClass {
	class := krt.FetchOne(ctx, gatewayClasses, krt.FilterKey(string(gc)))
	if class == nil {
		return &GatewayClass{
			Name:       string(gc),
			Controller: wellknown.GatewayControllerName, // TODO: make this configurable
		}
	}
	return class
}
