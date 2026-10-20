//go:build darwin

package vz

// ApplyStartOptionForTest runs a start option against fresh, empty start
// options and returns what it says; the test package cannot reach the
// unexported options type otherwise.
func ApplyStartOptionForTest(option VirtualMachineStartOption) error {
	return option(&virtualMachineStartOptions{})
}

// HostIsMacOS27ForTest is the verdict the provisioning option reads first.
func HostIsMacOS27ForTest() bool {
	return macOSAvailable(macOSGuestProvisioningVersion) == nil
}

// ProvisioningCompiledForTest is the verdict it reads second.
func ProvisioningCompiledForTest() bool {
	return provisioningCompiled()
}
