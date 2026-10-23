//go:build darwin

package osversion_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/forkcloser/vz/internal/osversion"
)

// TestBuildTargetKnowsEveryGatedVersion holds the version table to every
// macOS release an API in this module is gated on: a version missing from it
// is reported as "unsupported target version" whatever SDK the binary was
// built with, which reads like an SDK problem and is not one. Whether the
// build's SDK is new enough is the other outcome and depends on the machine,
// so only the table lookup is held here.
func TestBuildTargetKnowsEveryGatedVersion(t *testing.T) {
	t.Parallel()

	for _, version := range []float64{11, 12, 12.3, 13, 14, 15, 15.4, 26, 27} {
		err := osversion.MacOSBuildTargetAvailable(version)
		if err == nil {
			continue
		}

		if !errors.Is(err, osversion.ErrBuildTargetOSVersion) {
			t.Fatalf("version %v: %v, want ErrBuildTargetOSVersion or nil", version, err)
		}

		if strings.Contains(err.Error(), "unsupported target version") {
			t.Errorf("version %v is gated on by this module but missing from the build target table: %v", version, err)
		}
	}
}
