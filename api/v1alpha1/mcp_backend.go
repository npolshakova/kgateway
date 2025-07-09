package v1alpha1

// MCP configures mcp backends
type MCP struct {
	// Name is the backend name for this MCP configuration.
	// +required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Targets is a list of MCP targets to use for this backend.
	// +required
	// +kubebuilder:validation:MinItems=1
	Targets []McpTarget `json:"targets"`
}

// McpTarget defines a single MCP target configuration.
type McpTarget struct {
	// Name is the name of this MCP target.
	// +required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Host is the hostname or IP address of the MCP target.
	// +required
	// +kubebuilder:validation:MinLength=1
	Host string `json:"host"`

	// Port is the port number of the MCP target.
	// +required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`

	// EnableTls enables TLS for the connection to the MCP target.
	// +optional
	EnableTls bool `json:"enableTls,omitempty"`

	// Filters is a list of filters to apply to this target.
	// +optional
	Filters []McpFilter `json:"filters,omitempty"`
}

// McpFilter defines a filter configuration for MCP targets.
type McpFilter struct {
	// Type specifies the type of filter to apply.
	// +required
	// +kubebuilder:validation:MinLength=1
	Type string `json:"type"`

	// Match specifies the matching criteria for the filter.
	// +required
	Match FilterMatcher `json:"match"`
}

// FilterMatcher defines different matching strategies for filters.
type FilterMatcher struct {
	// Exact matches the exact string value.
	// +optional
	Exact *string `json:"exact,omitempty"`

	// Prefix matches strings that start with the specified prefix.
	// +optional
	Prefix *string `json:"prefix,omitempty"`

	// Suffix matches strings that end with the specified suffix.
	// +optional
	Suffix *string `json:"suffix,omitempty"`

	// Contains matches strings that contain the specified substring.
	// +optional
	Contains *string `json:"contains,omitempty"`

	// Regex matches strings using the specified regular expression.
	// +optional
	Regex *string `json:"regex,omitempty"`
}
