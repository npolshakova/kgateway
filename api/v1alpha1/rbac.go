package v1alpha1

// Rbac defines the configuration for role-based access control.
type Rbac struct {
	// Rules defines the RBAC rules for authorization.
	// +required
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	Rules []RbacRule `json:"rules"`

	// Action defines whether the rule allows or denies the request if matched.
	// If unspecified, the default is "Allow".
	// +kubebuilder:validation:Enum=Allow;Deny
	// +kubebuilder:default=Allow
	Action AuthorizationPolicyAction `json:"action,omitempty"`
}

// RbacRule defines a single RBAC rule.
type RbacRule struct {
	// Principal defines the authentication requirements for this rule.
	// +required
	Principal Principal `json:"principal"`

	// Access defines what resources and operations are allowed. If not specified, the rule will allow all resources and operations.
	Access AccessRule `json:"access,omitempty"`
}

// Principal defines authentication requirements that can be satisfied by different types of principals
type Principal struct {
	// JWTPrincipals defines JWT-based authentication rules.
	JWTPrincipals []JWTPrincipal `json:"jwt,omitempty"`
	// TODO: Add support for other principal types (e.g., CIDR ranges, header-based authentication, service accounts, IP ranges, etc.)
}

// AccessRule defines what resources and operations are allowed
type AccessRule struct {
	// Paths that have this prefix will be allowed.
	// +kubebuilder:validation:MinItems=1
	PathPrefixes []string `json:"pathPrefixes"`

	// What http methods (GET, POST, ...) are allowed.
	// +kubebuilder:validation:MinItems=1
	Methods []string `json:"methods"`
}

type AuthorizationPolicyAction string

const (
	AuthorizationPolicyActionAllow AuthorizationPolicyAction = "Allow"
	AuthorizationPolicyActionDeny  AuthorizationPolicyAction = "Deny"
)

// Permission configures what permissions should be granted.
// If more than one field is added, all of them need to match.
type Permission struct {
	// Paths that have this prefix will be allowed.
	PathPrefix string `json:"pathPrefix,omitempty"`
	// What http methods (GET, POST, ...) are allowed.
	Methods []string `json:"methods,omitempty"`
}

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
