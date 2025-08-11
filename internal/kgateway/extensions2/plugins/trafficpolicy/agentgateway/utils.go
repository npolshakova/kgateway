package agentgateway

type ProviderNeededMap struct {
	// map filter_chain name -> provider name -> provider
	Providers map[string]map[string]*TrafficPolicyAgentGatewayExtensionIR
}

func (p *ProviderNeededMap) Add(filterChain, providerName string, provider *TrafficPolicyAgentGatewayExtensionIR) {
	if p.Providers == nil {
		p.Providers = make(map[string]map[string]*TrafficPolicyAgentGatewayExtensionIR)
	}
	if p.Providers[filterChain] == nil {
		p.Providers[filterChain] = make(map[string]*TrafficPolicyAgentGatewayExtensionIR)
	}
	p.Providers[filterChain][providerName] = provider
}
