package backendtlspolicy

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	envoycore "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	corev1 "k8s.io/api/core/v1"
)

// handles conversion into envoy auth types
// based on https://github.com/solo-io/gloo/blob/main/projects/gloo/pkg/utils/ssl.go#L76

var noKeyFoundMsg = "no key ca.crt found"

func ResolveUpstreamSslConfig(cm *corev1.ConfigMap, sni string) (*envoyauth.UpstreamTlsContext, error) {
	common, err := ResolveCommonSslConfig(cm, false)
	if err != nil {
		return nil, err
	}

	return &envoyauth.UpstreamTlsContext{
		CommonTlsContext: common,
		Sni:              sni,
	}, nil
}

func ResolveCommonSslConfig(cm *corev1.ConfigMap, mustHaveCert bool) (*envoyauth.CommonTlsContext, error) {
	caCrt, err := getSslSecrets(cm)
	if err != nil {
		return nil, err
	}

	if err = validateCACertificate(caCrt); err != nil {
		return nil, err
	}

	tlsContext := &envoyauth.CommonTlsContext{
		// default params
		TlsParams: &envoyauth.TlsParameters{},
	}

	validationCtx := &envoyauth.CommonTlsContext_ValidationContext{
		ValidationContext: &envoyauth.CertificateValidationContext{
			TrustedCa: &envoycore.DataSource{
				Specifier: &envoycore.DataSource_InlineString{
					InlineString: caCrt,
				},
			},
		},
	}
	// sanList := VerifySanListToMatchSanList(cs.GetVerifySubjectAltName())
	// if len(sanList) != 0 {
	// 	validationCtx.ValidationContext.MatchSubjectAltNames = sanList
	// }
	tlsContext.ValidationContextType = validationCtx
	return tlsContext, err
}

func getValidationContext(cm *corev1.ConfigMap) (*envoyauth.CertificateValidationContext, error) {
	caCrt, err := getSslSecrets(cm)
	if err != nil {
		return nil, err
	}
	if err = validateCACertificate(caCrt); err != nil {
		return nil, err
	}
	return &envoyauth.CertificateValidationContext{
		TrustedCa: &envoycore.DataSource{
			Specifier: &envoycore.DataSource_InlineString{
				InlineString: caCrt,
			},
		},
	}, nil
}

func validateCACertificate(pemData string) error {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return errors.New("failed to decode PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	if cert.NotAfter.Before(time.Now()) {
		return errors.New("CA certificate has expired")
	}

	// TODO: add other validation (ex. IsCA?)

	return nil
}

func getSslSecrets(cm *corev1.ConfigMap) (string, error) {
	caCrt, ok := cm.Data["ca.crt"]
	if !ok {
		return "", errors.New(noKeyFoundMsg)
	}

	return caCrt, nil
}
