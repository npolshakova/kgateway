# Tracing e2e tests

The tracing e2e tests confirm the traces are present in the OTEL collector. The setup can be used with Jaeger to visualize the traces.

## Visualize traces with Jaeger

Make sure the `setup.yaml`, `tracing.yaml` and httpbin are applied.
```shell
kubectl apply -f testdata/setup.yaml
kubectl apply -f testdata/tracing.yaml
kubectl apply -f test/e2e/defaults/testdata/httpbin.yaml
```

The otel collector already has an exporter configured to:
```yaml
      zipkin:
        endpoint: http://tracing.default.svc.cluster.local:9411/api/v2/spans
```

Apply the jaeger config to setup the endpoint:
```shell
kubectl apply -f testdata/jaeger.yaml
```

Expose the gateway (ie. `kubectl port-forward svc/gw 8080:808`) and send some example requests:
```shell
curl localhost:8080/get -H "Host: www.example.com"
```

Then you can view the traces in Jaeger UI:
```shell
kubectl port-forward svc/tracing 16686:80
```