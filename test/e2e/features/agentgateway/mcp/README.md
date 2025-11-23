# MCP Tests details

For the dynamic routing we use two MCP servers one "user" and other "admin". 

Probably it can be simplified. For now we use two different docker images.

- User container uses `ghcr.io/peterj/mcp-website-fetcher:main` as a source and is copied to the CI registry on GitHub as `ghcr.io/kgateway-dev/mcp-website-fetcher:0.0.1`
- Admin container is built using upstream and placed `ghcr.io/kgateway-dev/mcp-admin-server:0.0.1` it can be rebuilt using Dockerfile in this directory. Use `docker build -t ghcr.io/kgateway-dev/mcp-admin-server:<version> .` and then push it to GH with `docker push ghcr.io/kgateway-dev/mcp-admin-server:<version>` command

## Keycloak MCP Authentication

For testing MCP Authentication with keycloak, apply the `remote-authn-keycloak.yaml`, `common.yaml` and `keycloak.yaml` 
manifests along with the `curl_pod.yaml`.

Send a request without authorization header:
```shell
curl -v -X POST http://gw.default.svc.cluster.local:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json,text/event-stream" \
  -d '{
    "method": "tools/call",
    "params": {
      "name": "fetch",
      "arguments": { "url": "http://google.com" }
    }
  }'
```

You should get a 401 Unauthorized response:
```shell
< HTTP/1.1 401 Unauthorized
< www-authenticate: Bearer resource_metadata="http://mcp-website-fetcher.default.svc.cluster.local/.well-known/oauth-protected-resource/mcp"
< content-type: application/json
< content-length: 65
< date: Thu, 20 Nov 2025 14:32:26 GMT
< 
* Connection #0 to host localhost left intact
{"error":"unauthorized","error_description":"JWT token required"}%
```

From the `curl` pod you should be able to get a token from keycloak under `access_token` and save it as `TOKEN`:
```shell
curl -s -X POST "http://keycloak.default:7080/realms/mcp/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials" \
  -d "client_id=mcp_proxy" \
  -d "client_secret=supersecret"
```

Then use the token to make the request:
```shell
curl -v -X POST http://gw.default.svc.cluster.local:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json,text/event-stream" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "method": "tools/call",
    "params": {
      "name": "fetch",
      "arguments": { "url": "http://google.com" }
    }
  }'
```

You can also test this with the MCP inspector:
1. Run `npx modelcontextprotocol/inspector#0.16.2` 
2. Open the inspector UI 
3. Attempt to connect to the port-forwarded gateway (`http://localhost:8080/`) without the token using Streamable HTTP
4. Set the `TOKEN` under the API Token Authentication field, then click Connect
5. Go to the `tools` tab and test the `fetch` tool with a URL of your choice
