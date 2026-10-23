//go:build darwin

package vz_test

import (
	"errors"
	"testing"

	"github.com/forkcloser/vz"
)

// TestProvisioningGateOrder holds the order the option refuses in: the host's
// version first, whatever the SDK; the SDK second, only on a macOS 27 host.
// The two facts the option reads are read here too, so the one outcome this
// host must produce is known before the option runs. A build that checks them
// the other way round fails where both gates fire: a host older than 27 built
// against an SDK older than 27, the macos-26 runner among them.
func TestProvisioningGateOrder(t *testing.T) {
	t.Parallel()

	// Complete values: on a macOS 27 host with the macOS 27 SDK the option
	// reaches the framework, which validates them (an empty full name is
	// refused), and that refusal is not the gate under test.
	option := vz.WithMacOSGuestProvisioning(vz.MacGuestProvisioningOptions{
		FullName: "Zero Gravity",
		Username: "zg",
		Password: "secret",
	})

	err := vz.ApplyStartOptionForTest(option)

	switch {
	case !vz.HostIsMacOS27ForTest():
		if !errors.Is(err, vz.ErrUnsupportedOSVersion) {
			t.Fatalf(
				"on a host older than macOS 27 the option = %v, want ErrUnsupportedOSVersion whatever the SDK",
				err,
			)
		}
	case !vz.ProvisioningCompiledForTest():
		if !errors.Is(err, vz.ErrProvisioningSDK) {
			t.Fatalf("on a macOS 27 host built against an older SDK the option = %v, want ErrProvisioningSDK", err)
		}
	default:
		if err != nil {
			t.Fatalf("on a macOS 27 host with the macOS 27 SDK the option = %v, want nil", err)
		}
	}
}
