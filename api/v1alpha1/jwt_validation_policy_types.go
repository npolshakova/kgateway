package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// +kubebuilder:rbac:groups=gateway.kgateway.dev,resources=jwtvalidationpolicies,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.kgateway.dev,resources=jwtvalidationpolicies/status,verbs=get;update;patch

// JWTValidationPolicy defines how JSON Web Tokens (JWTs) are extracted from requests and validated.
// Requests with invalid JWTs will be rejected. Requests without JWTs are allowed to proceed but
// will not have an associated authenticated identity. To enforce access only for authenticated
// users, this policy should be used in conjunction with an authorization policy.
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:metadata:labels={app=kgateway,app.kubernetes.io/name=kgateway}
// +kubebuilder:resource:categories=kgateway
// +kubebuilder:subresource:status
// +kubebuilder:metadata:labels="gateway.networking.k8s.io/policy=Direct"
type JWTValidationPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec JWTValidationPolicySpec `json:"spec,omitempty"`

	Status gwv1alpha2.PolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type JWTValidationPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JWTValidationPolicy `json:"items"`
}

type JWTValidationPolicySpec struct {
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	TargetRefs []LocalPolicyTargetReference `json:"targetRefs,omitempty"`

	// Map of JWT provider name to Provider.
	// If multiple providers are specified, the providers will be `OR`-ed together and will allow validation to any of the providers.
	// +kubebuilder:validation:MinProperties=1
	Providers map[string]*JWTProvider `json:"providers"`

	// TODO: add ValidationMode (REQUIRE_VALID,ALLOW_MISSING,ALLOW_MISSING_OR_FAILED). Check if supported by agentgateway
}

type JWTProvider struct {
	// Issuer of the JWT. the 'iss' claim of the JWT must match this.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Issuer string `json:"issuer"`

	// Audience is the list of audiences to be used for the JWT provider.
	// An incoming JWT must have an 'aud' claim, and it must be in this list.
	// +kubebuilder:validation:MinItems=1
	Audiences []string `json:"audiences"`

	// TokenSource is the token source to be used for the JWT provider.
	TokenSource *JWTTokenSource `json:"tokenSource,omitempty"`

	// ClaimToHeaders is the list of claims to headers to be used for the JWT provider.
	// Optionally set the claims from the JWT payload that you want to extract and add as headers
	// to the request before the request is forwarded to the upstream destination.
	ClaimToHeaders []*JWTClaimToHeader `json:"claimToHeaders,omitempty"`

	// JWKS is the source for the JSON Web Keys to be used to validate the JWT.
	JWKS JWKS `json:"jwks"`
}

type JWTTokenSource struct {
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Name string `json:"name"`
}

// JWTClaimToHeader allows copying verified claims to headers sent upstream
type JWTClaimToHeader struct {
	// Name is the JWT claim name, for example, "sub".
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Name string `json:"name"`

	// Header is the header the claim will be copied to, for example, "x-sub".
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Header string `json:"header"`
}

// JWKKind represents JWK type
// +kubebuilder:validation:Enum=REMOTE;LOCAL
type JWKKind string

const (
	Remote JWKKind = "REMOTE" // Equal
	Local  JWKKind = "LOCAL"
)

type JWKS struct {
	// +kubebuilder:validation:Required
	Kind JWKKind `json:"kind"`

	// LocalJwks configures provide a PEM-formatted public key or file to verify the JWT token.
	LocalJwks *LocalJWKS `json:"localJwks,omitempty"`

	// RemoteJwks configures a reference to the JSON Web Key Set (JWKS) server
	RemoteJwks *RemoteJWKS `json:"remoteJwks,omitempty"`
}

// LocalJWKSKind represents local JWKS type
// +kubebuilder:validation:Enum=INLINE;FILE;SECRET
type LocalJWKSKind string

const (
	JWKInline LocalJWKSKind = "INLINE" // Equal
	JWKFile   LocalJWKSKind = "FILE"
	JWKSecret LocalJWKSKind = "SECRET"
)

// LocalJWKS configures getting the public keys to validate the JWT from a local source, such as a Kubernetes secret,
// inline, raw string JWKS or file source.
type LocalJWKS struct {
	// +kubebuilder:validation:Required
	Kind LocalJWKSKind `json:"kind"`

	// File is the path to the file containing the JWKS
	File *string `json:"file,omitempty"`

	// InlineKey is the JWKS key as the raw, inline JWKS string
	InlineKey *string `json:"key,omitempty"`

	// SecretRef configures storing the JWK in a Kubernetes secret in the same namespace as the JWTValidationPolicy.
	SecretRef *corev1.LocalObjectReference `json:"secretRef,omitempty"`
}

// RemoteJWKS configures getting the public keys from a remote JSON Web Key Set (JWKS) server.
// This server must be accessible from your cluster.
type RemoteJWKS struct {
	// TargetRef is the target reference to the JWKS server.
	// If the JWKS server runs in your cluster, the destination can be a Kubernetes Service or kgateway Backend.
	// If the JWKS server runs outside your cluster, the destination should be a static kgateway Backend.
	// +kubebuilder:validation:Required
	TargetRef LocalPolicyTargetReference `json:"targetRefs"`

	// URL is used when accessing the upstream for Json Web Key Set.
	// This is used to set the host and path in the request
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Format=uri
	URL string `json:"url"`
}
