package jwt

import (
	"net/http"
	"path/filepath"

	"github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kgateway-dev/kgateway/v2/pkg/utils/fsutils"
	"github.com/kgateway-dev/kgateway/v2/test/gomega/matchers"
)

var (
	// manifests
	setupManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "setup.yaml")
	jwtManifest   = filepath.Join(fsutils.MustGetThisDir(), "testdata", "jwt.yaml")

	// Core infrastructure objects that we need to track
	gatewayObjectMeta = metav1.ObjectMeta{
		Name:      "gw",
		Namespace: "default",
	}
	gatewayService    = &corev1.Service{ObjectMeta: gatewayObjectMeta}
	gatewayDeployment = &appsv1.Deployment{ObjectMeta: gatewayObjectMeta}

	httpbinObjectMeta = metav1.ObjectMeta{
		Name:      "httpbin",
		Namespace: "httpbin",
	}
	httpbinDeployment = &appsv1.Deployment{ObjectMeta: httpbinObjectMeta}

	// Matches
	expectedJwtMissingFailedResponse = &matchers.HttpResponse{
		StatusCode: http.StatusUnauthorized,
		Body:       gomega.ContainSubstring("Jwt is missing"),
	}
	expectedJwtVerificationFailedResponse = &matchers.HttpResponse{
		StatusCode: http.StatusUnauthorized,
		Body:       gomega.ContainSubstring("Jwt verification fails"),
	}

	expectStatus200Success = &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Body:       nil,
	}

	// invalid jwt (not signed with correct key)
	badJwtToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2Rldi5leGFtcGxlLmNvbSIsImV4cCI6NDgwNDMyNDczNiwiaWF0IjoxNjQ4NjUxMTM2LCJvcmciOiJpbnRlcm5hbCIsImVtYWlsIjoiZGV2MkBrZ2F0ZXdheS5pbyIsImdyb3VwIjoiZW5naW5lZXJpbmciLCJzY29wZSI6ImlzOmRldmVsb3BlciJ9.pduAl6C0YofLSTUNcQuSd5dvrN-B8eE0pbOJJ9h5Fyh-k1HQQzSpZ47HJngclFmfcWk25qyJfLOnuVuA4PV6PwanPovL5YpdLlAbjHZPfDwsR1v8zUzb97yl-hbQzYCiA8coHO6rQE8hOYD59-DXkH6acuU8nVm3sv6VUA8zR5XpxZfJHJfRu8TZUFowk3FFrdh3nUSeeXLtm0YxN9uVEHKe3v_UEdMBUzri7wC1saKy7CcpikpBwd7itPMpT87BL_f1LvJf7LUEChRC-sp2LYsyjT-rme4YufPp1vVi5dMSCpfmvB1XlgFKzmGBPKvDJPta1DNOmHqEmKmgOQBCmw"

	/*
		Configured with these fields:
			{
			  "iss": "https://dev.example.com",
			  "exp": 4804324736,
			  "iat": 1648651136,
			  "org": "internal",
			  "email": "dev1@kgateway.io",
			  "group": "engineering",
			  "scope": "is:developer"
			}
		Using https://jwt.io/ and the following instructions to generate a public/private key pair:
		1. openssl genrsa 2048 > private-key.pem
		2. openssl rsa -in private-key.pem -pubout
		3. cat private-key.pem | pbcopy
	*/
	// claim has email=dev1@kgateway.io
	dev1JwtToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2Rldi5leGFtcGxlLmNvbSIsImV4cCI6NDgwNDMyNDczNiwiaWF0IjoxNjQ4NjUxMTM2LCJvcmciOiJpbnRlcm5hbCIsImVtYWlsIjoiZGV2MUBrZ2F0ZXdheS5pbyIsImdyb3VwIjoiZW5naW5lZXJpbmciLCJzY29wZSI6ImlzOmRldmVsb3BlciJ9.pqzk87Gny6mT8Gk7CVfkminm3u9CrNPhRt0oElwmfwZ7Jak1Ss4iOZ7MSZEgZFPxGiaz3DQyvos65dqbM_e4RaLYXb9fFYylaBl8kE8bhqMnXfPBNp9C4XTsSz4mR-eUvnkXXZ31dhMkoZvwIswWXR50wZ0rC6NF60Tye0sHJRdDcwL5778wDzLnualvtIiL-CbhWzXgRmjcrK3sbikLCHBjQiTEyBMPOVqS5NqJBgd7ZW1UASoxuxjCLsN8tBIaAFSACf8FZggAh9vEUJ_uc39kvOKQ0vs0pxvoYtsMPcndBYhws6IUhx_iF__qs_zz9mDNp8aMbXSlEdJG30wiRA"
)
