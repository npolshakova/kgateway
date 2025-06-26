package agentgatewaysyncer

const (
	TargetTypeResourceUrl = "type.googleapis.com/istio.adp.Resource"
	TargetTypeAddressUrl  = "type.googleapis.com/istio.workload.Address"

	MCPProtocol = "kgateway.dev/mcp"
	A2AProtocol = "kgateway.dev/a2a"

	MCPPathAnnotation = "kgateway.dev/mcp-path"
	A2APathAnnotation = "kgateway.dev/a2a-path"
)
