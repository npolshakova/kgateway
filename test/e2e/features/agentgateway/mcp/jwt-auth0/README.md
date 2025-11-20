Based on agentgateway auth0 example: 

1. Build docker 

```shell
docker build -t auth0-mock .
```

```shell
kind load --name kind docker-image auth0-mock:latest
```

```shell
curl -s -X POST http://auth0-mock.default:9000/register \
  -H "Content-Type: application/json" \
  -d '{
    "client_name": "Test Client",
    "redirect_uris": ["http://localhost/callback"]
  }'
```

```shell
curl -s -X POST http://auth0-mock.default:9000/register \
  -H "Content-Type: application/json" \
  -d '{
        "client_name": "Test Client",
        "redirect_uris": ["http://localhost/callback"]
      }'

```

Then copy the client_id from the response and use it in /authorize:

```shell
export CLIENT_ID=mcp_1cOsuWFas7Ii0IPvU4Y_OTdolYtueHDz

curl -s "http://auth0-mock.default:9000/authorize?response_type=code&client_id=CLIENT_ID&redirect_uri=http://localhost/callback&scope=openid%20profile&prompt=consent"
```

```shell
curl http://auth0-mock.default:9000/authorize

CODE=$(curl -s "http://auth0-mock.default:9000/authorize?response_type=code&client_id=$CLIENT_ID&redirect_uri=$REDIRECT_URI&scope=openid%20profile&prompt=consent" \
  | grep -o 'window.location.href *= *[^;]*' \
  | sed -E "s/window.location.href *= *'[^?]*\?code=([^']+)'.*/\1/")

echo $CODE

export CODE=GWIB3oLBJ5fyD4TzMklCF9_YKbsxXHgugd8k5OaAxE4
export CLIENT_ID=mcp_1cOsuWFas7Ii0IPvU4Y_OTdolYtueHDz
export CLIENT_SECRET=secret_DohY5CqgH0gczhyQKToOdxxc1ulWciZ8
export REDIRECT_URI=http://0.0.0.0/callback

curl -s -X POST http://auth0-mock.default:9000/token \
  -u "$CLIENT_ID:$CLIENT_SECRET" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code" \
  -d "code=$CODE" \
  -d "redirect_uri=$REDIRECT_URI"

```


***

```shell
curl -s "http://auth0-mock.default:9000/authorize?response_type=code&client_id=mcp_gi3APARn2_uHv2oxfJJqq2yZBDV4OyNo&redirect_uri=http://127.0.0.1/callback&scope=openid%20profile&prompt=consent"
```


```shell
Private JWK:
 {"d":"KEhvRCLz5YUve9YwaXQs8zJD6vvRygiaT60pkyADzAINVRIOsezVOUCX-aFzS14e5ioBKC9pYj74mPMoro1QlbmZMMbsps_9xO9iKSS3M83S7NckFOoZuJCmikRvQqahGbrokvUuEmS81ydMg6t5tCF918oBfL6A72DklnNgRxXDqv1ohY5h_z7z7eKjBXcxsF1S1Aybhbh4bOtiCgEb5Y5tlmIY_nDEsEw1oYAr2qj1Yib-Wbv5Lh05__kFvVM7DNKb6pq9PmZezMBzDP58BhMhsjt3rbdClE7gjD3ooAOGaMlRoZkP580SIC1hNVfrYibMBdH3GNze4q9cwR40KQ","dp":"5tQTAkVdUgFH4bxP6Nwit9b5mgdM112GwkUJnVoTctAADh7OPB3C18-6GqGE50-qv1AJWwEsP73t8CqrIyOFEs4iB7kKuucKZetnkbvjnbxFSHdSBrnlDkBkD4ajPXDx-KQ4HIjAnMvl0_sWDWrOd93G7aQFSHw2KbJnJMXx9d0","dq":"oDD86ei2GpPglWMvboNXHRdrtBNKDHRRLXI82ux1lvxD9wQbLBLhNLz3FgNFEiaKburNcExpPqRegEBifAiZvbHSOCb915MJS-x_MHob_TgHFA-IPIZvdFTSnFXBncPjIbT-J2ow5o3VR5_tE-kpmXfxl3GFcl6I4fUP-kxJjUU","e":"AQAB","kty":"RSA","n":"rmQv7yCC0-yR4zIirW1o04DrWRay04pXpugM7Kw7fdqiLIyFwxA6aQ_f1whyDvS3e-kGJnJ-byis9aNf9DJG_bYzdw99JFODlogY-AZL6449AeujbsaF_kxepZZiwkh9wLQ8zthr4ccVR7HM1AYlE___4ulNYwT5h5xuJOGZRsCb30I2AmDklFlPmtRGC5p7Pz5ZM0XBSfwRBPPa4mbWrghN_HhXrSDyj-eB1sp82clDvbI4hxet9_dW60jqTyj7xlh8jtl6yCviWNH_rqr8N3gvC1zFVYUJ6SbsY1LCHDWfXdqk8iYjtShTws4n5JX7VSHjqL1B5rdAo166eiqSGQ","p":"6hfY4tlo9rZ6roy1eF5T_YTYOCEU2gf0fsPkhF_KbdUoEgcn4CJZdqKsEgftz_dC5Ivfgy1cBhv3a377K5UW7RbDA1oGhTLPgwgdbhkc6b1_2bQFhb3mCY2EA7wET3OKenTrvX-xjccl3e1lsCoiJnduvVxJDQj7r4WMH-rrjeU","q":"vrYRgh8Eq5vrYG8vhWQd5RE-JBd8MBPbAnUUdmfFfs-neLq7rojUgMU7Z3WLQOwmcA1zW7zfvFp7yYxDoyx4gwvafV8FipAVM3SreWLddVWZ8VN_1c4s7Sv7tVp7jxhgfnROUI5NxHpvXYVCszC_dcKQVhpMKe2NkxGXmHIx0CU","qi":"PMk6mvu0hWMCD42vYg_gFG9jvztSMAS9O6JSw4d_ZAGPydhh0pTu38GrqVMFRH_7V389QpmgFJhYJ4QwqLH3zeTyG5OPB-zfoBJz2QSpLvGTTkgzr306pA57lhoMKN2adnF7HXBI2QRC1aInCqH40m8iQlTX-xR9LaWxasfhV9Q"}
```

```shell
Public JWK:
 {"e":"AQAB","kty":"RSA","n":"rmQv7yCC0-yR4zIirW1o04DrWRay04pXpugM7Kw7fdqiLIyFwxA6aQ_f1whyDvS3e-kGJnJ-byis9aNf9DJG_bYzdw99JFODlogY-AZL6449AeujbsaF_kxepZZiwkh9wLQ8zthr4ccVR7HM1AYlE___4ulNYwT5h5xuJOGZRsCb30I2AmDklFlPmtRGC5p7Pz5ZM0XBSfwRBPPa4mbWrghN_HhXrSDyj-eB1sp82clDvbI4hxet9_dW60jqTyj7xlh8jtl6yCviWNH_rqr8N3gvC1zFVYUJ6SbsY1LCHDWfXdqk8iYjtShTws4n5JX7VSHjqL1B5rdAo166eiqSGQ"}
```

```shell
JWKS:
 {
  "keys": [
    {
      "kid": "5333780687551038659",
      "e": "AQAB",
      "kty": "RSA",
      "n": "rmQv7yCC0-yR4zIirW1o04DrWRay04pXpugM7Kw7fdqiLIyFwxA6aQ_f1whyDvS3e-kGJnJ-byis9aNf9DJG_bYzdw99JFODlogY-AZL6449AeujbsaF_kxepZZiwkh9wLQ8zthr4ccVR7HM1AYlE___4ulNYwT5h5xuJOGZRsCb30I2AmDklFlPmtRGC5p7Pz5ZM0XBSfwRBPPa4mbWrghN_HhXrSDyj-eB1sp82clDvbI4hxet9_dW60jqTyj7xlh8jtl6yCviWNH_rqr8N3gvC1zFVYUJ6SbsY1LCHDWfXdqk8iYjtShTws4n5JX7VSHjqL1B5rdAo166eiqSGQ"
    }
  ]
}
```

```shell
JWT:
 eyJhbGciOiJSUzI1NiIsImtpZCI6IjUzMzM3ODA2ODc1NTEwMzg2NTkifQ.eyJhdWQiOiJhY2NvdW50IiwiZXhwIjoxNzYzNjc2Nzc2LCJpYXQiOjE3NjM2NzMxNzYsImlzcyI6Imh0dHBzOi8va2dhdGV3YXkuZGV2Iiwic3ViIjoidXNlckBrZ2F0ZXdheS5kZXYifQ.Fko5TMFRRJoXyidRaAmzmwlVHIwNxCXqiKf5BRw_sumTnpNmt9Qt_2RUQCn7tTC_gAV50FyV4WKwoyTzAn0S8mmgZumI8E2-Uoq-A8wAohz9rt4a61_gaDeXXn0dF3YitQicR30Q_buoi2Nki6ZRPf9FyE5ulO4Ut_PyQrNXwlwO7vr_U3DXfrzvT9y2aDdNndPr1GB4fWTM84mEdQgx3XevIc7yjnbgKHnvIRp4gEyh-QL0ZYisjD-tZIDloZoSZjNFYu6PIdoxAaz9WhINAkAqX9KS8cd6uO36nPDoDOT1UmCT2VBjNszhLaZqtRKbJUb1HYrn-Gzq8vumLn8sjQ
```
