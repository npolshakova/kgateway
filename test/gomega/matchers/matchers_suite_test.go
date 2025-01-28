//go:build ignore

package matchers_test

import (
	"testing"

	"github.com/solo-io/go-utils/testutils"

	. "github.com/onsi/ginkgo/v2"
)

func TestMatchers(t *testing.T) {
	testutils.RegisterCommonFailHandlers()
	RunSpecs(t, "Matchers Suite")
}
