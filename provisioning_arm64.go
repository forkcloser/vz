//go:build darwin

package vz

/*
#cgo darwin CFLAGS: -mmacosx-version-min=11 -x objective-c -fno-objc-arc
#cgo darwin LDFLAGS: -lobjc -framework Foundation -framework Virtualization
# include "virtualization_13_arm64.h"
# include "virtualization_27_arm64.h"
// The 13 header gives the start options, the 27 header the provisioning.
*/
import "C"

import (
	"errors"
)

// FORK: the macOS 27 guest provisioning API, which upstream does not bind yet.

// ErrProvisioningSDK reports a binary built against an SDK older than macOS
// 27: the provisioning calls were compiled out, so the option refuses instead
// of raising at run time.
var ErrProvisioningSDK = errors.New("guest provisioning needs a build against the macOS 27 SDK or newer")

// macOSGuestProvisioningVersion is the first macOS whose framework provisions
// a guest on first boot.
const macOSGuestProvisioningVersion = 27

// MacGuestProvisioningOptions is what a macOS 27 guest's Setup Assistant
// consumes on the first boot after a restore, so that the boot is unattended:
// the account it creates and two switches.
//
// macOS evaluates the options on that first boot only and a guest older than
// macOS 27 ignores them.
type MacGuestProvisioningOptions struct {
	// FullName is the account's display name.
	FullName string
	// Username is the account's short name, the one ssh logs in as.
	Username string
	// Password is the account's password.
	Password string
	// LogsInAutomatically logs the account in at startup.
	LogsInAutomatically bool
	// EnablesRemoteLogin turns Remote Login (ssh) on.
	EnablesRemoteLogin bool
}

// WithMacOSGuestProvisioning provisions a macOS guest on its first boot after
// a restore. It needs macOS 27 on the host and a binary built against the
// macOS 27 SDK; otherwise Start returns an error before touching the VM.
//
// It builds the macOS start options when no earlier option did, so it
// combines with WithStartUpFromMacOSRecovery only when that one comes first.
func WithMacOSGuestProvisioning(provisioning MacGuestProvisioningOptions) VirtualMachineStartOption {
	return func(vmso *virtualMachineStartOptions) error {
		if err := macOSAvailable(macOSGuestProvisioningVersion); err != nil {
			return err
		}

		if !provisioningCompiled() {
			return ErrProvisioningSDK
		}

		if vmso.macOSVirtualMachineStartOptionsPtr == nil {
			vmso.macOSVirtualMachineStartOptionsPtr = C.newVZMacOSVirtualMachineStartOptions(C.bool(false))
		}

		fullName := charWithGoString(provisioning.FullName)
		defer fullName.Free()

		username := charWithGoString(provisioning.Username)
		defer username.Free()

		password := charWithGoString(provisioning.Password)
		defer password.Free()

		nserrPtr := newNSErrorAsNil()
		C.setGuestProvisioningOptionsVZMacOSVirtualMachineStartOptions(
			vmso.macOSVirtualMachineStartOptionsPtr,
			fullName.CString(),
			username.CString(),
			password.CString(),
			C.bool(provisioning.LogsInAutomatically),
			C.bool(provisioning.EnablesRemoteLogin),
			&nserrPtr,
		)

		// newNSError returns *NSError: a nil pointer in an error is not nil.
		//revive:disable-next-line:if-return
		if err := newNSError(nserrPtr); err != nil {
			return err
		}

		return nil
	}
}

// provisioningCompiled reports whether the SDK this binary was built against
// had the provisioning calls; the C side knows, the Go side asks.
func provisioningCompiled() bool {
	return bool(C.vzMacGuestProvisioningCompiled())
}
