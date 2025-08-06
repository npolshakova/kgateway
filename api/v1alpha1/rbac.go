package v1alpha1

// Rbac defines the configuration for role-based access control.
type Rbac struct {
	// Policies defines a list of roles and the principals that are assigned/denied the role.
	// A policy matches if and only if at least one of its permissions match the action taking place
	// AND at least one of its principals match the downstream
	// AND the condition is true if specified.
	// +required
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	Policies []RbacPolicy `json:"policies"`

	// Action defines whether the rule allows or denies the request if matched.
	// If unspecified, the default is "Allow".
	// +kubebuilder:validation:Enum=Allow;Deny
	// +kubebuilder:default=Allow
	Action AuthorizationPolicyAction `json:"action,omitempty"`
}

// RbacPolicy defines a single RBAC rule.
type RbacPolicy struct {
	// Principals defines the list of authentication requirements for this rule.
	// +optional
	Principals []Principal `json:"principals,omitempty"`

	//// Permissions defines what resources and operations are allowed. If not specified, the rule will allow all resources and operations.
	//Permissions []Permission `json:"permissions,omitempty"`

	// Condition defines a set of conditions that must be satisfied for the rule to match.
	Conditions *CELConditions `json:"conditions,omitempty"`
}

// Principal defines authentication requirements that can be satisfied by different types of principals
type Principal struct {
	// JWTPrincipals defines a map of provider name to JWT principals (e.g. "my-provider") used in RBAC
	// These must match provider names defined in the GatewayExtension and configured in the JWT TrafficPolicy field.
	// +optional
	JWTPrincipals map[string]*JWTPrincipal `json:"jwt,omitempty"`

	// TODO: support other principal types (e.g. OIDC, etc.)
}

type CELConditions struct {
	// CelMatchExpression defines a set of conditions that must be satisfied for the rule to match.
	// +kubebuilder:validation:MinItems=1
	CelMatchExpression []string `json:"matchExpressions,omitempty"`
}

type AuthorizationPolicyAction string

const (
	AuthorizationPolicyActionAllow AuthorizationPolicyAction = "Allow"
	AuthorizationPolicyActionDeny  AuthorizationPolicyAction = "Deny"
)

// JWTPrincipal defines JWT-based authentication configuration
type JWTPrincipal struct {
	// Claims defines a set of claims that make up the principal.
	// All listed claims must be present and match the given value (AND semantics).
	// Commonly, the 'iss' and 'sub' or 'email' claims are used.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	Claims []JWTClaimMatch `json:"claims"`
}

// JWTClaimMatch configures the claim to match
type JWTClaimMatch struct {
	// Name is the name of the claim to match (e.g., "sub", "role"). It can be a nested claim of type
	// (eg. "my.cool.key") where the nested claim delimiter must use dot "." to separate the name path.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	Name string `json:"name"`

	// Value is the expected value of the claim.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	Value string `json:"value"`

	// Matcher specifies how the claim value should be matched.
	// +kubebuilder:validation:Enum=Exact;Boolean;Contains
	// +kubebuilder:default=Exact
	Matcher ClaimMatcher `json:"matcher,omitempty"`
}

// ClaimMatcher specifies how claims should be matched to the value.
type ClaimMatcher string

const (
	// ClaimMatcherExactString indicates the claim value is a string that exactly matches the value.
	ClaimMatcherExactString ClaimMatcher = "Exact"
	// ClaimMatcherBoolean indicates the claim value is a boolean that matches the value.
	ClaimMatcherBoolean ClaimMatcher = "Boolean"
	// ClaimMatcherContains indicates the claim value is a list that contains a string that matches the value.
	ClaimMatcherContains ClaimMatcher = "Contains"
)
