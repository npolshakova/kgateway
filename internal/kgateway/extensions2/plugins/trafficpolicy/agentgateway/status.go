package agentgateway

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	extensionsplug "github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugin"
)

func getPolicyStatusFn(
	cl client.Client,
) extensionsplug.GetPolicyStatusFn {
	return func(ctx context.Context, nn types.NamespacedName) (gwv1alpha2.PolicyStatus, error) {
		res := v1alpha1.TrafficPolicy{}
		err := cl.Get(ctx, nn, &res)
		if err != nil {
			return gwv1alpha2.PolicyStatus{}, err
		}
		return res.Status, nil
	}
}

func patchPolicyStatusFn(
	cl client.Client,
) extensionsplug.PatchPolicyStatusFn {
	return func(ctx context.Context, nn types.NamespacedName, policyStatus gwv1alpha2.PolicyStatus) error {
		res := v1alpha1.TrafficPolicy{}
		err := cl.Get(ctx, nn, &res)
		if err != nil {
			return err
		}

		res.Status = policyStatus
		if err := cl.Status().Patch(ctx, &res, client.Merge); err != nil {
			return fmt.Errorf("error updating status for TrafficPolicy %s: %w", nn.String(), err)
		}
		return nil
	}
}
