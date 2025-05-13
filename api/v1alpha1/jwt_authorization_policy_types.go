package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// +kubebuilder:rbac:groups=gateway.kgateway.dev,resources=jwtauthorizationpolicies,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.kgateway.dev,resources=jwtauthorizationpolicies/status,verbs=get;update;patch

// JWTAuthorizationPolicy defines rules for how JSON Web Tokens (JWTs) are used to define authorization rules.
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:metadata:labels={app=kgateway,app.kubernetes.io/name=kgateway}
// +kubebuilder:resource:categories=kgateway
// +kubebuilder:subresource:status
// +kubebuilder:metadata:labels="gateway.networking.k8s.io/policy=Direct"
type JWTAuthorizationPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec JWTAuthorizationPolicySpec `json:"spec,omitempty"`

	Status gwv1alpha2.PolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type JWTAuthorizationPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JWTAuthorizationPolicy `json:"items"`
}

type JWTAuthorizationPolicySpec struct {
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	TargetRefs []LocalPolicyTargetReference `json:"targetRefs,omitempty"`

	// JWT Principals are the list of JWT principals to be used for the JWT provider.
	// TODO: Rename to Rules? This is really a list of authorization rules based on principal claims?
	// +kubebuilder:validation:MinItems=1
	Principals []JWTPrincipals `json:"principals"`

	// Action defines whether the rule allows or denies the request if matched.
	// If unspecified, the default is "Allow".
	// +kubebuilder:validation:Enum=ALLOW;DENY
	Action *JWTAuthorizationPolicyAction `json:"action,omitempty"`
}

type JWTAuthorizationPolicyAction string

const (
	JWTAuthorizationPolicyActionAllow JWTAuthorizationPolicyAction = "ALLOW"
	JWTAuthorizationPolicyActionDeny  JWTAuthorizationPolicyAction = "DENY"
)

type JWTPrincipals struct {
	// TODO: is there any benefit of having this as a map?
	// RequiredClaims defines a set of claims that make up the principal.
	// All listed claims must be present and match the given value (AND semantics).
	// Commonly, the 'iss' and 'sub' or 'email' claims are used.
	// +kubebuilder:validation:MinItems=1
	RequiredClaims []JWTClaimMatch `json:"requiredClaims"`

	// Optional: Allow if the JWT is issued by one of these providers (by name).
	// If empty, all configured providers are allowed.
	AllowedProviders []string `json:"allowedProviders,omitempty"`
}

// TODO: support nested claims
type JWTClaimMatch struct {
	// Name is the name of the claim to match (e.g., "sub", "role").
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Value is the expected value of the claim.
	// +kubebuilder:validation:MinLength=1
	Value string `json:"value"`

	// TODO: ClaimMatcher type? (exact string,bool,list contains, etc.)
}
