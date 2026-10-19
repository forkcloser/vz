//
//  virtualization_27_arm64.h
//
//  FORK: guest provisioning for macOS guests (macOS 27 API).
//

#pragma once

#import "internal/osversion/virtualization_helper.h"
#import <Virtualization/Virtualization.h>

/* macOS 27 API */
// Whether this binary was compiled against an SDK that has
// VZMacGuestProvisioningOptions. The Go side refuses the option otherwise,
// instead of raising at run time.
bool vzMacGuestProvisioningCompiled(void);
// Attaches freshly built provisioning options to a VZMacOSVirtualMachineStartOptions.
// Every string is copied; error receives the framework's NSError on refusal.
void setGuestProvisioningOptionsVZMacOSVirtualMachineStartOptions(void *startOptions,
    const char *fullName, const char *username, const char *password,
    bool logsInAutomatically, bool enablesRemoteLogin, void **error);
