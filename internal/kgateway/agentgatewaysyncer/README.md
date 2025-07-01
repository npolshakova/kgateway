# agentgateway syncer

This syncer configures xds updates for the [agentgateway](https://agentgateway.dev/) data plane. 

### APIs

The syncer uses the following APIs:

- [workload](https://github.com/agentgateway/agentgateway/tree/main/go/api/workload.pb.go)
- [resource](https://github.com/agentgateway/agentgateway/tree/main/go/api/resource.pb.go)

### Conformance tests

```
make run HELM_ADDITIONAL_VALUES=test/kubernetes/e2e/tests/manifests/agent-gateway-integration.yaml; CONFORMANCE_GATEWAY_CLASS=agentgateway make conformance 
```