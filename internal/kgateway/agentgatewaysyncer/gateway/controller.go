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
	"strconv"

	"go.uber.org/atomic"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayalpha "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gateway "sigs.k8s.io/gateway-api/apis/v1beta1"

	networkingclient "istio.io/client-go/pkg/apis/networking/v1"
	kubesecrets "istio.io/istio/pilot/pkg/credentials/kube"
	"istio.io/istio/pilot/pkg/features"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/serviceregistry/kube/controller"
	"istio.io/istio/pilot/pkg/status"
	"istio.io/istio/pkg/cluster"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"
	"istio.io/istio/pkg/config/schema/gvk"
	"istio.io/istio/pkg/config/schema/gvr"
	"istio.io/istio/pkg/config/schema/kind"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/kube/kubetypes"
	istiolog "istio.io/istio/pkg/log"
	"istio.io/istio/pkg/revisions"
	"istio.io/istio/pkg/slices"
	"istio.io/istio/pkg/util/sets"
	"istio.io/istio/pkg/workloadapi"
)

var log = istiolog.RegisterScope("gateway", "gateway-api controller")

var errUnsupportedOp = fmt.Errorf("unsupported operation: the gateway config store is a read-only view")

// Controller defines the controller for the gateway-api. The controller reads a variety of resources (Gateway types, as well
// as adjacent types like Namespace and Service), and through `krt`, translates them into Istio types (Gateway/VirtualService).
//
// Most resources are fully "self-contained" with krt, but there are a few usages breaking out of `krt`; these are managed by `krt.RecomputeProtected`.
// These are recomputed on each new PushContext initialization, which will call Controller.Reconcile().
//
// The generated Istio types are not stored in the cluster at all and are purely internal. Calls to List() (from PushContext)
// will expose these. They can be introspected at /debug/configz.
//
// The status on all gateway-api types is also tracked. Each collection emits downstream objects, but also status about the
// input type. If the status changes, it is queued to asynchronously update the status of the object in Kubernetes.
type Controller struct {
	// client for accessing Kubernetes
	client kube.Client

	// the cluster where the gateway-api controller runs
	cluster cluster.ID
	// revision the controller is running under
	revision string

	// status controls the status writing queue. Status will only be written if statusEnabled is true, which
	// is only the case when we are the leader.
	status *StatusCollections

	waitForCRD func(class schema.GroupVersionResource, stop <-chan struct{}) bool

	// gatewayContext exposes us to the internal Istio service registry. This is outside krt knowledge (currently), so,
	// so we wrap it in a RecomputeProtected.
	// Most usages in the API are directly referenced typed objects (Service, ServiceEntry, etc) so this is not needed typically.
	gatewayContext krt.RecomputeProtected[*atomic.Pointer[GatewayContext]]
	// tagWatcher allows us to check which tags are ours. Unlike most Istio codepaths, we read istio.io/rev=<tag> and not just
	// revisions for Gateways. This is because a Gateway is sort of a mix of a Deployment and Config.
	// Since the TagWatcher is not yet krt-aware, we wrap this in RecomputeProtected.
	tagWatcher krt.RecomputeProtected[revisions.TagWatcher]

	stop chan struct{}

	xdsUpdater model.XDSUpdater

	// Handlers tracks all registered handlers, so that syncing can be detected
	handlers []krt.HandlerRegistration

	// outputs contains all the output collections for this controller.
	// Currently, the only usage of this controller is from non-krt things (PushContext) so this is not exposed directly.
	// If desired in the future, it could be.
	outputs Outputs
}

func (c *Controller) Collection() krt.Collection[model.ADPResource] {
	return c.outputs.ADPResources
}

type ParentInfo struct {
	Key  parentKey
	Info parentInfo
}

func (pi ParentInfo) ResourceName() string {
	return pi.Key.Name // TODO!!!! more infoi and section name
}

type TypedResource struct {
	Kind config.GroupVersionKind
	Name types.NamespacedName
}

type Outputs struct {
	Gateways        krt.Collection[Gateway]
	VirtualServices krt.Collection[*config.Config]
	ReferenceGrants ReferenceGrants

	ADPResources krt.Collection[model.ADPResource]
}

type Inputs struct {
	Namespaces krt.Collection[*corev1.Namespace]

	Services krt.Collection[*corev1.Service]
	Secrets  krt.Collection[*corev1.Secret]

	GatewayClasses  krt.Collection[*gateway.GatewayClass]
	Gateways        krt.Collection[*gateway.Gateway]
	HTTPRoutes      krt.Collection[*gateway.HTTPRoute]
	GRPCRoutes      krt.Collection[*gatewayv1.GRPCRoute]
	TCPRoutes       krt.Collection[*gatewayalpha.TCPRoute]
	TLSRoutes       krt.Collection[*gatewayalpha.TLSRoute]
	ReferenceGrants krt.Collection[*gateway.ReferenceGrant]
	ServiceEntries  krt.Collection[*networkingclient.ServiceEntry]
	InferencePools  krt.Collection[*inf.InferencePool]
}

var _ model.GatewayController = &Controller{}

func NewController(
	kc kube.Client,
	waitForCRD func(class schema.GroupVersionResource, stop <-chan struct{}) bool,
	options controller.Options,
	xdsUpdater model.XDSUpdater,
) *Controller {
	stop := make(chan struct{})
	opts := krt.NewOptionsBuilder(stop, "gateway", options.KrtDebugger)

	tw := revisions.NewTagWatcher(kc, options.Revision)
	c := &Controller{
		client:         kc,
		cluster:        options.ClusterID,
		revision:       options.Revision,
		status:         &StatusCollections{},
		tagWatcher:     krt.NewRecomputeProtected(tw, false, opts.WithName("tagWatcher")...),
		waitForCRD:     waitForCRD,
		gatewayContext: krt.NewRecomputeProtected(atomic.NewPointer[GatewayContext](nil), false, opts.WithName("gatewayContext")...),
		stop:           stop,
		xdsUpdater:     xdsUpdater,
	}
	tw.AddHandler(func(s sets.String) {
		c.tagWatcher.TriggerRecomputation()
	})

	inputs := Inputs{
		Namespaces: krt.NewInformer[*corev1.Namespace](kc, opts.WithName("informer/Namespaces")...),
		Secrets: krt.WrapClient[*corev1.Secret](
			kclient.NewFiltered[*corev1.Secret](kc, kubetypes.Filter{
				FieldSelector: kubesecrets.SecretsFieldSelector,
				ObjectFilter:  kc.ObjectFilter(),
			}),
			opts.WithName("informer/Secrets")...,
		),
		Services: krt.WrapClient[*corev1.Service](
			kclient.NewFiltered[*corev1.Service](kc, kubetypes.Filter{ObjectFilter: kc.ObjectFilter()}),
			opts.WithName("informer/Services")...,
		),
		GatewayClasses: buildClient[*gateway.GatewayClass](c, kc, gvr.GatewayClass, opts, "informer/GatewayClasses"),
		Gateways:       buildClient[*gateway.Gateway](c, kc, gvr.KubernetesGateway, opts, "informer/Gateways"),
		HTTPRoutes:     buildClient[*gateway.HTTPRoute](c, kc, gvr.HTTPRoute, opts, "informer/HTTPRoutes"),
		GRPCRoutes:     buildClient[*gatewayv1.GRPCRoute](c, kc, gvr.GRPCRoute, opts, "informer/GRPCRoutes"),

		ReferenceGrants: buildClient[*gateway.ReferenceGrant](c, kc, gvr.ReferenceGrant, opts, "informer/ReferenceGrants"),
		ServiceEntries:  buildClient[*networkingclient.ServiceEntry](c, kc, gvr.ServiceEntry, opts, "informer/ServiceEntries"),
		InferencePools:  buildClient[*inf.InferencePool](c, kc, gvr.InferencePool, opts, "informer/InferencePools"),
	}
	if features.EnableAlphaGatewayAPI {
		inputs.TCPRoutes = buildClient[*gatewayalpha.TCPRoute](c, kc, gvr.TCPRoute, opts, "informer/TCPRoutes")
		inputs.TLSRoutes = buildClient[*gatewayalpha.TLSRoute](c, kc, gvr.TLSRoute, opts, "informer/TLSRoutes")
	} else {
		// If disabled, still build a collection but make it always empty
		inputs.TCPRoutes = krt.NewStaticCollection[*gatewayalpha.TCPRoute](nil, nil, opts.WithName("disable/TCPRoutes")...)
		inputs.TLSRoutes = krt.NewStaticCollection[*gatewayalpha.TLSRoute](nil, nil, opts.WithName("disable/TLSRoutes")...)
	}

	handlers := []krt.HandlerRegistration{}

	GatewayClassStatus, GatewayClasses := GatewayClassesCollection(inputs.GatewayClasses, opts)
	registerStatus(c, GatewayClassStatus)

	ReferenceGrants := BuildReferenceGrants(ReferenceGrantsCollection(inputs.ReferenceGrants, opts))

	// GatewaysStatus cannot is not fully complete until its join with route attachments to report attachedRoutes.
	// Do not register yet.
	GatewaysStatus, Gateways := GatewayCollection(
		inputs.Gateways,
		GatewayClasses,
		inputs.Namespaces,
		ReferenceGrants,
		inputs.Secrets,
		options.DomainSuffix,
		c.gatewayContext,
		c.tagWatcher,
		opts,
	)
	ports := krt.NewIndex(Gateways, func(o Gateway) []string {
		return []string{fmt.Sprint(o.parentInfo.Port)}
	}).AsCollection(opts.WithName("PortBindings")...)
	Binds := krt.NewManyCollection(ports, func(ctx krt.HandlerContext, object krt.IndexObject[string, Gateway]) []model.ADPResource {
		port, _ := strconv.Atoi(object.Key)
		uniq := sets.New[types.NamespacedName]()
		for _, gw := range object.Objects {
			uniq.Insert(types.NamespacedName{
				Namespace: gw.parent.Namespace,
				Name:      gw.parent.Name,
			})
		}
		return slices.Map(uniq.UnsortedList(), func(e types.NamespacedName) model.ADPResource {
			bind := Bind{
				Bind: &workloadapi.Bind{
					Key:  object.Key + "/" + e.String(),
					Port: uint32(port),
				},
			}
			return toResource(e, bind)
		})
	}, opts.WithName("Binds")...)
	WaypointBinds := krt.NewCollection(inputs.Gateways, func(ctx krt.HandlerContext, gw *gateway.Gateway) *model.ADPResource {
		if gw.Spec.GatewayClassName != "istio-waypoint" {
			return nil
		}
		port := 15008
		e := config.NamespacedName(gw)
		bind := Bind{
			Bind: &workloadapi.Bind{
				Key:  "waypoint/" + e.String(),
				Port: uint32(port),
			},
		}
		return toResourcep(e, bind)
	}, opts.WithName("WaypointBinds")...)
	Listeners := krt.NewCollection(Gateways, func(ctx krt.HandlerContext, obj Gateway) *model.ADPResource {
		l := &workloadapi.Listener{
			Key:         obj.ResourceName(),
			Name:        string(obj.parentInfo.SectionName),
			BindKey:     fmt.Sprint(obj.parentInfo.Port) + "/" + obj.parent.Namespace + "/" + obj.parent.Name,
			GatewayName: obj.parent.Namespace + "/" + obj.parent.Name,
			Hostname:    obj.parentInfo.OriginalHostname,
		}

		switch obj.parentInfo.Protocol {
		case gatewayv1.HTTPProtocolType:
			l.Protocol = workloadapi.Protocol_HTTP
		case gatewayv1.HTTPSProtocolType:
			l.Protocol = workloadapi.Protocol_HTTPS
			if obj.TLSInfo == nil {
				return nil
			}
			l.Tls = &workloadapi.TLSConfig{
				Cert:       obj.TLSInfo.Cert,
				PrivateKey: obj.TLSInfo.Key,
			}
		case gatewayv1.TLSProtocolType:
			l.Protocol = workloadapi.Protocol_TLS
			if obj.TLSInfo == nil {
				return nil
			}
			l.Tls = &workloadapi.TLSConfig{
				Cert:       obj.TLSInfo.Cert,
				PrivateKey: obj.TLSInfo.Key,
			}
		case gatewayv1.TCPProtocolType:
			l.Protocol = workloadapi.Protocol_TCP
		default:
			return nil
		}
		return toResourcep(types.NamespacedName{
			Namespace: obj.parent.Namespace,
			Name:      obj.parent.Name,
		}, ADPListener{l})
	}, opts.WithName("Listeners")...)
	WaypointListeners := krt.NewCollection(inputs.Gateways, func(ctx krt.HandlerContext, gw *gateway.Gateway) *model.ADPResource {
		if gw.Spec.GatewayClassName != "istio-waypoint" {
			return nil
		}

		e := config.NamespacedName(gw)
		bind := ADPListener{
			Listener: &workloadapi.Listener{
				Key:         "waypoint/" + e.String(),
				Name:        "waypoint/" + e.String(),
				BindKey:     "waypoint/" + e.String(),
				GatewayName: e.String(),
				Hostname:    "",
				Protocol:    workloadapi.Protocol_HBONE,
				Tls:         nil,
			},
		}
		return toResourcep(e, bind)
	}, opts.WithName("WaypointListeners")...)

	RouteParents := BuildRouteParents(Gateways)

	routeInputs := RouteContextInputs{
		Grants:          ReferenceGrants,
		RouteParents:    RouteParents,
		DomainSuffix:    options.DomainSuffix,
		Services:        inputs.Services,
		Namespaces:      inputs.Namespaces,
		ServiceEntries:  inputs.ServiceEntries,
		InferencePools:  inputs.InferencePools,
		internalContext: c.gatewayContext,
	}
	ADPRoutes := ADPRouteCollection(
		inputs.HTTPRoutes,
		routeInputs,
		opts,
	)
	//tcpRoutes := TCPRouteCollection(
	//	inputs.TCPRoutes,
	//	routeInputs,
	//	opts,
	//)
	//registerStatus(c, tcpRoutes.Status)
	//tlsRoutes := TLSRouteCollection(
	//	inputs.TLSRoutes,
	//	routeInputs,
	//	opts,
	//)
	//registerStatus(c, tlsRoutes.Status)
	httpRoutes := HTTPRouteCollection(
		inputs.HTTPRoutes,
		routeInputs,
		opts,
	)
	registerStatus(c, httpRoutes.Status)
	status, _ := krt.NewStatusCollection(inputs.InferencePools, func(krtctx krt.HandlerContext, obj *inf.InferencePool) (
		*inf.InferencePoolStatus,
		*any,
	) {
		status := obj.Status.DeepCopy()
		myGws := sets.New[types.NamespacedName]()
		allGws := sets.New[types.NamespacedName]() // this is dumb but https://github.com/kubernetes-sigs/gateway-api-inference-extension/issues/942...
		allGwsRaw := krt.Fetch(krtctx, inputs.Gateways)
		for _, g := range allGwsRaw {
			allGws.Insert(config.NamespacedName(g))
			if string(g.Spec.GatewayClassName) == features.GatewayAPIDefaultGatewayClass {
				myGws.Insert(config.NamespacedName(g))
			}
		}
		seen := sets.New[types.NamespacedName]()
		np := []inf.PoolStatus{}
		for _, s := range status.Parents {
			k := types.NamespacedName{
				Name:      s.GatewayRef.Name,
				Namespace: s.GatewayRef.Namespace,
			}
			if !allGws.Contains(k) {
				// Even if it's not ours, delete stale ref. Shrug.
				continue
			}
			if s.GatewayRef.Kind != gvk.KubernetesGateway.Kind {
				np = append(np, s)
				continue
			}
			if seen.Contains(k) {
				continue
			}
			if !myGws.Contains(k) {
				np = append(np, s)
				continue
			}
			myGws.Delete(k)
			seen.Insert(k)
			conds := map[string]*condition{
				string(inf.InferencePoolConditionAccepted): {
					reason:  string(inf.InferencePoolReasonAccepted),
					message: "Referenced by an HTTPRoute accepted by the parentRef Gateway",
				},
			}
			np = append(np, inf.PoolStatus{
				GatewayRef: corev1.ObjectReference{
					APIVersion: gatewayv1.GroupVersion.String(),
					Kind:       gvk.KubernetesGateway.Kind,
					Namespace:  k.Namespace,
					Name:       k.Name,
				},
				Conditions: setConditions(obj.Generation, s.Conditions, conds),
			})
		}
		for _, k := range myGws.UnsortedList() {
			conds := map[string]*condition{
				string(inf.InferencePoolConditionAccepted): {
					reason:  string(inf.InferencePoolReasonAccepted),
					message: "Referenced by an HTTPRoute accepted by the parentRef Gateway",
				},
			}
			np = append(np, inf.PoolStatus{
				GatewayRef: corev1.ObjectReference{
					APIVersion: gatewayv1.GroupVersion.String(),
					Kind:       gvk.KubernetesGateway.Kind,
					Namespace:  k.Namespace,
					Name:       k.Name,
				},
				Conditions: setConditions(obj.Generation, nil, conds),
			})
		}
		status.Parents = np
		return status, nil
	}, opts.WithName("InferencePools")...)
	registerStatus(c, status)
	//grpcRoutes := GRPCRouteCollection(
	//	inputs.GRPCRoutes,
	//	routeInputs,
	//	opts,
	//)
	//registerStatus(c, grpcRoutes.Status)

	RouteAttachments := krt.JoinCollection([]krt.Collection[*RouteAttachment]{
		// tcpRoutes.RouteAttachments,
		// tlsRoutes.RouteAttachments,
		httpRoutes.RouteAttachments,
		// grpcRoutes.RouteAttachments,
	}, opts.WithName("RouteAttachments")...)
	RouteAttachmentsIndex := krt.NewIndex(RouteAttachments, func(o *RouteAttachment) []GatewayAndListener {
		return []GatewayAndListener{{
			ListenerName: o.ListenerName,
			To:           o.To,
		}}
	})

	GatewayFinalStatus := FinalGatewayStatusCollection(GatewaysStatus, RouteAttachments, RouteAttachmentsIndex, opts)
	registerStatus(c, GatewayFinalStatus)

	VirtualServices := krt.JoinCollection([]krt.Collection[*config.Config]{
		// tcpRoutes.VirtualServices,
		// tlsRoutes.VirtualServices,
		httpRoutes.VirtualServices,
		// grpcRoutes.VirtualServices,
	}, opts.WithName("DerivedVirtualServices")...)

	ADPResources := krt.JoinCollection([]krt.Collection[model.ADPResource]{Binds, WaypointBinds, Listeners, WaypointListeners, ADPRoutes}, opts.WithName("ADPResources")...)

	outputs := Outputs{
		ReferenceGrants: ReferenceGrants,
		Gateways:        Gateways,
		VirtualServices: VirtualServices,

		ADPResources: ADPResources,
	}
	c.outputs = outputs

	handlers = append(handlers,
		outputs.VirtualServices.RegisterBatch(pushXds(xdsUpdater,
			func(t *config.Config) model.ConfigKey {
				return model.ConfigKey{
					Kind:      kind.VirtualService,
					Name:      t.Name,
					Namespace: t.Namespace,
				}
			}), false),
		outputs.ADPResources.RegisterBatch(pushXds(xdsUpdater,
			func(t model.ADPResource) model.ConfigKey {
				return model.ConfigKey{
					Kind: kind.ADP,
					Name: t.ResourceName(),
				}
			}), false),
		outputs.Gateways.RegisterBatch(pushXds(xdsUpdater,
			func(t Gateway) model.ConfigKey {
				return model.ConfigKey{
					Kind:      kind.Gateway,
					Name:      t.Name,
					Namespace: t.Namespace,
				}
			}), false),
		outputs.ReferenceGrants.collection.RegisterBatch(pushXds(xdsUpdater,
			func(t ReferenceGrant) model.ConfigKey {
				return model.ConfigKey{
					Kind:      kind.KubernetesGateway,
					Name:      t.Source.Name,
					Namespace: t.Source.Namespace,
				}
			}), false))
	c.handlers = handlers

	return c
}

// buildClient is a small wrapper to build a krt collection based on a delayed informer.
func buildClient[I controllers.ComparableObject](
	c *Controller,
	kc kube.Client,
	res schema.GroupVersionResource,
	opts krt.OptionsBuilder,
	name string,
) krt.Collection[I] {
	filter := kclient.Filter{
		ObjectFilter: kubetypes.ComposeFilters(kc.ObjectFilter(), c.inRevision),
	}

	// all other types are filtered by revision, but for gateways we need to select tags as well
	if res == gvr.KubernetesGateway {
		filter.ObjectFilter = kc.ObjectFilter()
	}

	cc := kclient.NewDelayedInformer[I](kc, res, kubetypes.StandardInformer, filter)
	return krt.WrapClient[I](cc, opts.WithName(name)...)
}

func (c *Controller) Schemas() collection.Schemas {
	return collection.SchemasFor(
		collections.VirtualService,
		collections.Gateway,
	)
}

func (c *Controller) Get(typ config.GroupVersionKind, name, namespace string) *config.Config {
	return nil
}

func (c *Controller) List(typ config.GroupVersionKind, namespace string) []config.Config {
	switch typ {
	case gvk.Gateway:
		res := slices.MapFilter(c.outputs.Gateways.List(), func(g Gateway) *config.Config {
			if g.Valid {
				return g.Config
			}
			return nil
		})
		return res
	case gvk.VirtualService:
		return slices.Map(c.outputs.VirtualServices.List(), func(e *config.Config) config.Config {
			return *e
		})
	default:
		return nil
	}
}

func (c *Controller) SetStatusWrite(enabled bool, statusManager *status.Manager) {
	if enabled && features.EnableGatewayAPIStatus && statusManager != nil {
		var q status.Queue = statusManager.CreateGenericController(func(status status.Manipulator, context any) {
			status.SetInner(context)
		})
		c.status.SetQueue(q)
	} else {
		c.status.UnsetQueue()
	}
}

// Reconcile is called each time the `gatewayContext` may change. We use this to mark it as updated.
func (c *Controller) Reconcile(ps *model.PushContext) {
	ctx := NewGatewayContext(ps, c.cluster)
	c.gatewayContext.Modify(func(i **atomic.Pointer[GatewayContext]) {
		(*i).Store(&ctx)
	})
	c.gatewayContext.MarkSynced()
}

func (c *Controller) Create(config config.Config) (revision string, err error) {
	return "", errUnsupportedOp
}

func (c *Controller) Update(config config.Config) (newRevision string, err error) {
	return "", errUnsupportedOp
}

func (c *Controller) UpdateStatus(config config.Config) (newRevision string, err error) {
	return "", errUnsupportedOp
}

func (c *Controller) Patch(orig config.Config, patchFn config.PatchFunc) (string, error) {
	return "", errUnsupportedOp
}

func (c *Controller) Delete(typ config.GroupVersionKind, name, namespace string, _ *string) error {
	return errUnsupportedOp
}

func (c *Controller) RegisterEventHandler(typ config.GroupVersionKind, handler model.EventHandler) {
}

func (c *Controller) Run(stop <-chan struct{}) {
	if features.EnableGatewayAPIGatewayClassController {
		go func() {
			if c.waitForCRD(gvr.GatewayClass, stop) {
				gcc := NewClassController(c.client)
				c.client.RunAndWait(stop)
				gcc.Run(stop)
			}
		}()
	}

	tw := c.tagWatcher.AccessUnprotected()
	go tw.Run(stop)
	go func() {
		kube.WaitForCacheSync("gateway tag watcher", stop, tw.HasSynced)
		c.tagWatcher.MarkSynced()
	}()

	<-stop
	close(c.stop)
}

func (c *Controller) HasSynced() bool {
	if !(c.outputs.VirtualServices.HasSynced() &&
		c.outputs.Gateways.HasSynced() &&
		c.outputs.ReferenceGrants.collection.HasSynced()) {
		return false
	}
	for _, h := range c.handlers {
		if !h.HasSynced() {
			return false
		}
	}
	return true
}

func (c *Controller) SecretAllowed(resourceName string, namespace string) bool {
	return c.outputs.ReferenceGrants.SecretAllowed(nil, resourceName, namespace)
}

func pushXds[T any](xds model.XDSUpdater, f func(T) model.ConfigKey) func(events []krt.Event[T]) {
	return func(events []krt.Event[T]) {
		if xds == nil {
			return
		}
		cu := sets.New[model.ConfigKey]()
		for _, e := range events {
			for _, i := range e.Items() {
				c := f(i)
				if c != (model.ConfigKey{}) {
					cu.Insert(c)
				}
			}
		}
		if len(cu) == 0 {
			return
		}
		xds.ConfigUpdate(&model.PushRequest{
			Full:           true,
			ConfigsUpdated: cu,
			Reason:         model.NewReasonStats(model.ConfigUpdate),
		})
	}
}

func (c *Controller) inRevision(obj any) bool {
	object := controllers.ExtractObject(obj)
	if object == nil {
		return false
	}
	return config.LabelsInRevision(object.GetLabels(), c.revision)
}
