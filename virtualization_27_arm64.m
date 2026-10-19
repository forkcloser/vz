//
//  virtualization_27_arm64.m
//
//  FORK: guest provisioning for macOS guests (macOS 27 API).
//

#import "virtualization_27_arm64.h"

#ifdef __arm64__

bool vzMacGuestProvisioningCompiled(void)
{
#ifdef INCLUDE_TARGET_OSX_27
    return true;
#else
    return false;
#endif
}

void setGuestProvisioningOptionsVZMacOSVirtualMachineStartOptions(void *startOptions,
    const char *fullName, const char *username, const char *password,
    bool logsInAutomatically, bool enablesRemoteLogin, void **error)
{
#ifdef INCLUDE_TARGET_OSX_27
    if (@available(macOS 27, *)) {
        VZMacGuestProvisioningOptions *provisioning = [[VZMacGuestProvisioningOptions alloc] init];
        [provisioning setFullName:[NSString stringWithUTF8String:fullName]];
        [provisioning setUsername:[NSString stringWithUTF8String:username]];
        [provisioning setPassword:[NSString stringWithUTF8String:password]];
        [provisioning setLogsInAutomatically:(BOOL)logsInAutomatically];
        [provisioning setEnablesRemoteLogin:(BOOL)enablesRemoteLogin];
        // The property is `copy`: the start options keep their own copy, so
        // ours goes (no ARC in this binding).
        [(VZMacOSVirtualMachineStartOptions *)startOptions setGuestProvisioningOptions:provisioning
                                                                                 error:(NSError *_Nullable *_Nullable)error];
        [provisioning release];
        return;
    }
#endif
    RAISE_UNSUPPORTED_MACOS_EXCEPTION();
}

#endif
