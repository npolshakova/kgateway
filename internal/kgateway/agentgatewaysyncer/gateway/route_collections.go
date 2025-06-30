// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gateway

import (
	"fmt"
	"iter"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayalpha "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gateway "sigs.k8s.io/gateway-api/apis/v1beta1"

	istio "istio.io/api/networking/v1alpha3"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/constants"
	"istio.io/istio/pkg/config/schema/gvk"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
	"istio.io/istio/pkg/slices"
	"istio.io/istio/pkg/util/protomarshal"
)

func HTTPRouteCollection(
	httpRoutes krt.Collection[*gateway.HTTPRoute],
	inputs RouteContextInputs,
	opts krt.OptionsBuilder,
) RouteResult[*gateway.HTTPRoute, gateway.HTTPRouteStatus] {
	routeCount := gatewayRouteAttachmentCountCollection(inputs, httpRoutes, gvk.HTTPRoute, opts)
	baseVirtualServices := krt.NewManyCollection(httpRoutes, func(krtctx krt.HandlerContext, obj *gateway.HTTPRoute) []RouteWithKey {
		ctx := inputs.WithCtx(krtctx)
		route := obj.Spec
		parentRefs, gwResult := computeRoute(ctx, obj, func(obj *gateway.HTTPRoute) iter.Seq2[*istio.HTTPRoute, *ConfigError] {
			return func(yield func(*istio.HTTPRoute, *ConfigError) bool) {
				for n, r := range route.Rules {
					// split the rule to make sure each rule has up to one match
					matches := slices.Reference(r.Matches)
					if len(matches) == 0 {
						matches = append(matches, nil)
					}
					for _, m := range matches {
						if m != nil {
							r.Matches = []gateway.HTTPRouteMatch{*m}
						}
						if !yield(convertHTTPRoute(ctx, r, obj, n)) {
							return
						}
					}
				}
			}
		})

		count := 0
		virtualServices := []RouteWithKey{}
		for _, parent := range filteredReferences(parentRefs) {
			// for gateway routes, build one VS per gateway+host
			routeKey := parent.InternalName
			vsHosts := hostnameToStringList(route.Hostnames)
			routes := gwResult.routes
			if len(routes) == 0 {
				continue
			}
			// Create one VS per hostname with a single hostname.
			// This ensures we can treat each hostname independently, as the spec requires
			for _, h := range vsHosts {
				if !parent.hostnameAllowedByIsolation(h) {
					// TODO: standardize a status message for this upstream and report
					continue
				}
				name := fmt.Sprintf("%s-%d-%s", obj.Name, count, constants.KubernetesGatewayName)
				sortHTTPRoutes(routes)
				cfg := &Config{
					Meta: Meta{
						CreationTimestamp: obj.CreationTimestamp.Time,
						//GroupVersionKind:  gvk.VirtualService,
						Name:        name,
						Annotations: routeMeta(obj),
						Namespace:   obj.Namespace,
						Domain:      ctx.DomainSuffix,
					},
					Spec: &istio.VirtualService{
						Hosts:    []string{h},
						Gateways: []string{parent.InternalName},
						Http:     routes,
					},
				}
				virtualServices = append(virtualServices, RouteWithKey{
					Config: cfg,
					Key:    routeKey + "/" + h,
				})
				count++
			}
		}
		return virtualServices
	}, opts.WithName("HTTPRoute")...)

	finalVirtualServices := mergeHTTPRoutes(baseVirtualServices, opts.WithName("HTTPRouteMerged")...)
	return RouteResult[*gateway.HTTPRoute, gateway.HTTPRouteStatus]{
		VirtualServices:  finalVirtualServices,
		RouteAttachments: routeCount,
	}
}

func ADPRouteCollection(
	httpRoutes krt.Collection[*gateway.HTTPRoute],
	inputs RouteContextInputs,
	opts krt.OptionsBuilder,
) krt.Collection[ADPResource] {
	routes := krt.NewManyCollection(httpRoutes, func(krtctx krt.HandlerContext, obj *gateway.HTTPRoute) []ADPResource {
		ctx := inputs.WithCtx(krtctx)
		route := obj.Spec
		parentRefs, gwResult := computeRoute(ctx, obj, func(obj *gateway.HTTPRoute) iter.Seq2[ADPRoute, *ConfigError] {
			return func(yield func(ADPRoute, *ConfigError) bool) {
				for n, r := range route.Rules {
					// split the rule to make sure each rule has up to one match
					matches := slices.Reference(r.Matches)
					if len(matches) == 0 {
						matches = append(matches, nil)
					}
					for idx, m := range matches {
						if m != nil {
							r.Matches = []gateway.HTTPRouteMatch{*m}
						}
						res, err := convertHTTPRouteToADP(ctx, r, obj, n, idx)

						if !yield(ADPRoute{Route: res}, err) {
							return
						}
					}
				}
			}
		})

		var res []ADPResource
		for _, parent := range filteredReferences(parentRefs) {
			// for gateway routes, build one VS per gateway+host
			routes := gwResult.routes
			if len(routes) == 0 {
				continue
			}
			gw := types.NamespacedName{
				Namespace: parent.ParentKey.Namespace,
				Name:      parent.ParentKey.Name,
			}
			res = append(res, slices.Map(routes, func(e ADPRoute) ADPResource {
				inner := protomarshal.Clone(e.Route)
				_, name, _ := strings.Cut(parent.InternalName, "/")
				inner.ListenerKey = name
				inner.Key = inner.Key + "." + string(parent.ParentSection)
				return toResource(gw, ADPRoute{Route: inner})
			})...)
		}
		return res
	}, opts.WithName("ADPRoutes")...)

	return routes
}

type conversionResult[O any] struct {
	error  *ConfigError
	routes []O
}

func GRPCRouteCollection(
	grpcRoutes krt.Collection[*gatewayv1.GRPCRoute],
	inputs RouteContextInputs,
	opts krt.OptionsBuilder,
) RouteResult[*gatewayv1.GRPCRoute, gatewayv1.GRPCRouteStatus] {
	routeCount := gatewayRouteAttachmentCountCollection(inputs, grpcRoutes, gvk.GRPCRoute, opts)
	baseVirtualServices := krt.NewManyCollection(grpcRoutes, func(krtctx krt.HandlerContext, obj *gatewayv1.GRPCRoute) []RouteWithKey {
		ctx := inputs.WithCtx(krtctx)
		route := obj.Spec
		parentRefs, gwResult := computeRoute(ctx, obj, func(obj *gatewayv1.GRPCRoute) iter.Seq2[*istio.HTTPRoute, *ConfigError] {
			return func(yield func(*istio.HTTPRoute, *ConfigError) bool) {
				for n, r := range route.Rules {
					// split the rule to make sure each rule has up to one match
					matches := slices.Reference(r.Matches)
					if len(matches) == 0 {
						matches = append(matches, nil)
					}
					for _, m := range matches {
						if m != nil {
							r.Matches = []gatewayv1.GRPCRouteMatch{*m}
						}
						if !yield(convertGRPCRoute(ctx, r, obj, n)) {
							return
						}
					}
				}
			}
		})

		count := 0
		var virtualServices []RouteWithKey
		for _, parent := range filteredReferences(parentRefs) {
			// for gateway routes, build one VS per gateway+host
			routeKey := parent.InternalName
			vsHosts := hostnameToStringList(route.Hostnames)
			routes := gwResult.routes
			if len(routes) == 0 {
				continue
			}
			// Create one VS per hostname with a single hostname.
			// This ensures we can treat each hostname independently, as the spec requires
			for _, h := range vsHosts {
				if !parent.hostnameAllowedByIsolation(h) {
					// TODO: standardize a status message for this upstream and report
					continue
				}
				name := fmt.Sprintf("%s-%d-%s", obj.Name, count, constants.KubernetesGatewayName)
				sortHTTPRoutes(routes)
				cfg := &Config{
					Meta: Meta{
						CreationTimestamp: obj.CreationTimestamp.Time,
						//GroupVersionKind:  gvk.VirtualService,
						Name:        name,
						Annotations: routeMeta(obj),
						Namespace:   obj.Namespace,
						Domain:      ctx.DomainSuffix,
					},
					Spec: &istio.VirtualService{
						Hosts:    []string{h},
						Gateways: []string{parent.InternalName},
						Http:     routes,
					},
				}
				virtualServices = append(virtualServices, RouteWithKey{
					Config: cfg,
					Key:    routeKey + "/" + h,
				})
				count++
			}
		}
		return virtualServices
	}, opts.WithName("GRPCRoute")...)

	finalVirtualServices := mergeHTTPRoutes(baseVirtualServices, opts.WithName("GRPCRouteMerged")...)
	return RouteResult[*gatewayv1.GRPCRoute, gatewayv1.GRPCRouteStatus]{
		VirtualServices:  finalVirtualServices,
		RouteAttachments: routeCount,
	}
}

func TCPRouteCollection(
	tcpRoutes krt.Collection[*gatewayalpha.TCPRoute],
	inputs RouteContextInputs,
	opts krt.OptionsBuilder,
) RouteResult[*gatewayalpha.TCPRoute, gatewayalpha.TCPRouteStatus] {
	routeCount := gatewayRouteAttachmentCountCollection(inputs, tcpRoutes, gvk.TCPRoute, opts)
	virtualServices := krt.NewManyCollection(tcpRoutes, func(krtctx krt.HandlerContext, obj *gatewayalpha.TCPRoute) []*Config {
		ctx := inputs.WithCtx(krtctx)
		route := obj.Spec
		parentRefs, gwResult := computeRoute(ctx, obj,
			func(obj *gatewayalpha.TCPRoute) iter.Seq2[*istio.TCPRoute, *ConfigError] {
				return func(yield func(*istio.TCPRoute, *ConfigError) bool) {
					for _, r := range route.Rules {
						if !yield(convertTCPRoute(ctx, r, obj)) {
							return
						}
					}
				}
			})

		var vs []*Config
		for _, parent := range filteredReferences(parentRefs) {
			routes := gwResult.routes
			vsHosts := []string{"*"}
			for i, host := range vsHosts {
				name := fmt.Sprintf("%s-tcp-%d-%s", obj.Name, i, constants.KubernetesGatewayName)
				// Create one VS per hostname with a single hostname.
				// This ensures we can treat each hostname independently, as the spec requires
				vs = append(vs, &Config{
					Meta: Meta{
						CreationTimestamp: obj.CreationTimestamp.Time,
						//GroupVersionKind:  gvk.VirtualService,
						Name:        name,
						Annotations: routeMeta(obj),
						Namespace:   obj.Namespace,
						Domain:      ctx.DomainSuffix,
					},
					Spec: &istio.VirtualService{
						// We can use wildcard here since each listener can have at most one route bound to it, so we have
						// a single VS per Gateway.
						Hosts:    []string{host},
						Gateways: []string{parent.InternalName},
						Tcp:      routes,
					},
				})
			}
		}
		return vs
	}, opts.WithName("TCPRoute")...)

	return RouteResult[*gatewayalpha.TCPRoute, gatewayalpha.TCPRouteStatus]{
		VirtualServices:  virtualServices,
		RouteAttachments: routeCount,
	}
}

func TLSRouteCollection(
	tlsRoutes krt.Collection[*gatewayalpha.TLSRoute],
	inputs RouteContextInputs,
	opts krt.OptionsBuilder,
) RouteResult[*gatewayalpha.TLSRoute, gatewayalpha.TLSRouteStatus] {
	routeCount := gatewayRouteAttachmentCountCollection(inputs, tlsRoutes, gvk.TLSRoute, opts)
	virtualServices := krt.NewManyCollection(tlsRoutes, func(krtctx krt.HandlerContext, obj *gatewayalpha.TLSRoute) []*Config {
		ctx := inputs.WithCtx(krtctx)
		route := obj.Spec
		parentRefs, gwResult := computeRoute(ctx,
			obj, func(obj *gatewayalpha.TLSRoute) iter.Seq2[*istio.TLSRoute, *ConfigError] {
				return func(yield func(*istio.TLSRoute, *ConfigError) bool) {
					for _, r := range route.Rules {
						if !yield(convertTLSRoute(ctx, r, obj)) {
							return
						}
					}
				}
			})

		var vs []*Config
		for _, parent := range filteredReferences(parentRefs) {
			routes := gwResult.routes
			vsHosts := hostnameToStringList(route.Hostnames)
			for i, host := range vsHosts {
				name := fmt.Sprintf("%s-tls-%d-%s", obj.Name, i, constants.KubernetesGatewayName)
				filteredRoutes := routes
				// Create one VS per hostname with a single hostname.
				// This ensures we can treat each hostname independently, as the spec requires
				vs = append(vs, &Config{
					Meta: Meta{
						CreationTimestamp: obj.CreationTimestamp.Time,
						//GroupVersionKind:  gvk.VirtualService,
						Name:        name,
						Annotations: routeMeta(obj),
						Namespace:   obj.Namespace,
						Domain:      ctx.DomainSuffix,
					},
					Spec: &istio.VirtualService{
						Hosts:    []string{host},
						Gateways: []string{parent.InternalName},
						Tls:      filteredRoutes,
					},
				})
			}
		}
		return vs
	})
	return RouteResult[*gatewayalpha.TLSRoute, gatewayalpha.TLSRouteStatus]{
		VirtualServices:  virtualServices,
		RouteAttachments: routeCount,
	}
}

// computeRoute holds the common route building logic shared amongst all types
func computeRoute[T controllers.Object, O comparable](ctx RouteContext, obj T, translator func(
	obj T,
) iter.Seq2[O, *ConfigError],
) ([]routeParentReference, conversionResult[O]) {
	parentRefs := extractParentReferenceInfo(ctx, ctx.RouteParents, obj)

	convertRules := func() conversionResult[O] {
		res := conversionResult[O]{}
		for vs, err := range translator(obj) {
			// This was a hard error
			if vs == nil {
				res.error = err
				return conversionResult[O]{error: err}
			}
			// Got an error but also routes
			if err != nil {
				res.error = err
			}
			res.routes = append(res.routes, vs)
		}
		return res
	}
	gwResult := buildMeshAndGatewayRoutes(parentRefs, convertRules)

	return parentRefs, gwResult
}

// RouteContext defines a common set of inputs to a route collection. This should be built once per route translation and
// not shared outside of that.
// The embedded RouteContextInputs is typically based into a collection, then translated to a RouteContext with RouteContextInputs.WithCtx().
type RouteContext struct {
	Krt krt.HandlerContext
	RouteContextInputs
}

func (r RouteContext) LookupHostname(hostname string, namespace string) *model.Service {
	if c := r.internalContext.Get(r.Krt).Load(); c != nil {
		return c.GetService(hostname, namespace)
	}
	return nil
}

type RouteContextInputs struct {
	Grants         ReferenceGrants
	RouteParents   RouteParents
	DomainSuffix   string
	Services       krt.Collection[*corev1.Service]
	InferencePools krt.Collection[*inf.InferencePool]
	Namespaces     krt.Collection[*corev1.Namespace]
	ServiceEntries krt.Collection[*networkingclient.ServiceEntry]
}

func (i RouteContextInputs) WithCtx(krtctx krt.HandlerContext) RouteContext {
	return RouteContext{
		Krt:                krtctx,
		RouteContextInputs: i,
	}
}

type RouteWithKey struct {
	*Config
	Key string
}

func (r RouteWithKey) ResourceName() string {
	return config.NamespacedName(r.Config).String()
}

func (r RouteWithKey) Equals(o RouteWithKey) bool {
	return r.Config.Equals(o.Config)
}

// buildMeshAndGatewayRoutes contains common logic to build a set of routes with gateway semantics
func buildMeshAndGatewayRoutes[T any](parentRefs []routeParentReference, convertRules func() T) T {
	return convertRules()
}

// RouteResult holds the result of a route collection
type RouteResult[I, IStatus any] struct {
	// VirtualServices are the primary output that configures the internal routing logic
	VirtualServices krt.Collection[*Config]
	// RouteAttachments holds information about parent attachment to routes, used for computed the `attachedRoutes` count.
	RouteAttachments krt.Collection[*RouteAttachment]
}

type GatewayAndListener struct {
	// To is assumed to be a Gateway
	To           types.NamespacedName
	ListenerName string
}

func (g GatewayAndListener) String() string {
	return g.To.String() + "/" + g.ListenerName
}

type RouteAttachment struct {
	From TypedResource
	// To is assumed to be a Gateway
	To           types.NamespacedName
	ListenerName string
}

func (r *RouteAttachment) ResourceName() string {
	return r.From.Kind.String() + "/" + r.From.Name.String() + "/" + r.To.String() + "/" + r.ListenerName
}

func (r *RouteAttachment) Equals(other RouteAttachment) bool {
	return r.From == other.From && r.To == other.To && r.ListenerName == other.ListenerName
}

// gatewayRouteAttachmentCountCollection holds the generic logic to determine the parents a route is attached to, used for
// computing the aggregated `attachedRoutes` status in Gateway.
func gatewayRouteAttachmentCountCollection[T controllers.Object](
	inputs RouteContextInputs,
	col krt.Collection[T],
	kind config.GroupVersionKind,
	opts krt.OptionsBuilder,
) krt.Collection[*RouteAttachment] {
	return krt.NewManyCollection(col, func(krtctx krt.HandlerContext, obj T) []*RouteAttachment {
		ctx := inputs.WithCtx(krtctx)
		from := TypedResource{
			Kind: kind,
			Name: config.NamespacedName(obj),
		}

		parentRefs := extractParentReferenceInfo(ctx, inputs.RouteParents, obj)
		return slices.MapFilter(filteredReferences(parentRefs), func(e routeParentReference) **RouteAttachment {
			if e.ParentKey.Kind != gvk.KubernetesGateway {
				return nil
			}
			return ptr.Of(&RouteAttachment{
				From: from,
				To: types.NamespacedName{
					Name:      e.ParentKey.Name,
					Namespace: e.ParentKey.Namespace,
				},
				ListenerName: string(e.ParentSection),
			})
		})
	}, opts.WithName(kind.Kind+"/count")...)
}

// mergeHTTPRoutes merges HTTProutes by key. Gateway API has semantics for the ordering of `match` rules, that merges across resource.
// So we merge everything (by key) following that ordering logic, and sort into a linear list (how VirtualService semantics work).
func mergeHTTPRoutes(baseVirtualServices krt.Collection[RouteWithKey], opts ...krt.CollectionOption) krt.Collection[*Config] {
	groupedRoutes := krt.NewCollection(baseVirtualServices, func(ctx krt.HandlerContext, obj RouteWithKey) *IndexObject[string, RouteWithKey] {
		return &IndexObject[string, RouteWithKey]{
			Key:     obj.Key,
			Objects: []RouteWithKey{obj},
		}
	}, opts...)
	finalVirtualServices := krt.NewCollection(groupedRoutes, func(ctx krt.HandlerContext, object IndexObject[string, RouteWithKey]) **Config {
		configs := object.Objects
		if len(configs) == 1 {
			return &configs[0].Config
		}
		sortRoutesByCreationTime(configs)
		base := configs[0].DeepCopy()
		baseVS := base.Spec.(*istio.VirtualService)
		for _, config := range configs[1:] {
			thisVS := config.Spec.(*istio.VirtualService)
			baseVS.Http = append(baseVS.Http, thisVS.Http...)
			// append parents
			base.Annotations[constants.InternalParentNames] = fmt.Sprintf("%s,%s",
				base.Annotations[constants.InternalParentNames], config.Annotations[constants.InternalParentNames])
		}
		sortHTTPRoutes(baseVS.Http)
		return ptr.Of(&base)
	}, opts...)
	return finalVirtualServices
}
